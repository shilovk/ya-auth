package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	loginURL         = "https://login.yandex.ru/info?format=json"
	oauthStateString = "pseudo-random"
)

var (
	conf = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     endpoint,
	}
	endpoint = oauth2.Endpoint{
		AuthURL:  "https://oauth.yandex.ru/authorize",
		TokenURL: "https://oauth.yandex.ru/token",
	}
)

type User struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
	Email string `json:"default_email"`
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
		<html><body><a href="/login">Войти с помощью Яндекс</a></body></html>
	`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	token, err := conf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, loginURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "OAuth "+token.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var u User
	err = json.Unmarshal(content, &u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Пользователь: %+v\n", u)
	w.Write(content)
}
