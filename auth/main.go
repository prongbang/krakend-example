package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func token(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	data := `
	{
		"access_token": {
			"aud": "http://api.example.com",
			"iss": "https://myapp.example.com",
			"sub": "1234567890qwertyuio",
			"jti": "mnb23vcsrt756yuiomnbvcx98ertyuiop",
			"roles": ["admin", "user"],
			"exp": 1735689600
		},
		"refresh_token": {
			"aud": "http://api.example.com",
			"iss": "https://myapp.example.com",
			"sub": "1234567890qwertyuio",
			"jti": "mnb23vcsrt756yuiomn12876bvcx98ertyuiop",
			"exp": 1735689600
		},
		"exp": 1735689600
	}	
	`

	maps := map[string]interface{}{}
	json.Unmarshal([]byte(data), &maps)

	json.NewEncoder(w).Encode(maps)
}

func user(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	maps := map[string]interface{}{
		"id":   1,
		"name": "hello",
	}

	json.NewEncoder(w).Encode(maps)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	buf := new(bytes.Buffer)

	responseData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Expiration   int    `json:"exp"`
	}{}

	req, err := http.NewRequest("POST", "http://localhost:2222/refresh-token", new(bytes.Buffer))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &responseData)

	fmt.Println(buf.String())

	json.NewEncoder(w).Encode(responseData)
}

func main() {
	http.HandleFunc("/token", token)
	http.HandleFunc("/refresh-token", refreshToken)
	http.HandleFunc("/user", user)

	fmt.Println("Start server: http://localhost:8800/")

	http.ListenAndServe(":8800", nil)
}
