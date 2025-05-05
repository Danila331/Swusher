package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var (
	clientID     = "87f7692d6fc44163887fcef0998dbd0a"
	clientSecret = "8ef462fd8bbb4b54bf9cd88441feb333"
	redirectURI  = "http://localhost:8080/oauth/yandex/callback"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&redirect_uri=%s",
		clientID, url.QueryEscape(redirectURI),
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Нет кода", http.StatusBadRequest)
		return
	}

	// Получаем access_token
	resp, err := http.PostForm("https://oauth.yandex.ru/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	})
	if err != nil {
		http.Error(w, "Ошибка получения токена", 500)
		log.Println("Token error:", err)
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Ошибка парсинга токена", 500)
		return
	}

	// Получаем данные пользователя
	req, _ := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	req.Header.Set("Authorization", "OAuth "+tokenResp.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка запроса данных", 500)
		return
	}
	defer userResp.Body.Close()

	var userData map[string]interface{}
	json.NewDecoder(userResp.Body).Decode(&userData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

func main() {
	http.HandleFunc("/oauth/yandex/login", loginHandler)
	http.HandleFunc("/oauth/yandex/callback", callbackHandler)

	fmt.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
