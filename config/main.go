package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := `
	{
		"keys": [
			{
				"kty": "oct",
				"alg": "A128KW",
				"k": "GawgguFyGrWKav7AX4VKUg",
				"kid": "sim1"
			},
			{
				"kty": "oct",
				"k": "AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow",
				"kid": "sim2",
				"alg": "HS256"
			}
		]
	}
	`

	maps := map[string]interface{}{}
	json.Unmarshal([]byte(data), &maps)

	json.NewEncoder(w).Encode(maps)
}

func main() {
	http.HandleFunc("/jwk/symmetric.json", greet)

	fmt.Println("Start server: http://localhost:5555/")

	http.ListenAndServe(":5555", nil)
}
