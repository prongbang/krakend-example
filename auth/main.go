package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	krakendjose "github.com/devopsfaith/krakend-jose"
	ginkrakendjose "github.com/devopsfaith/krakend-jose/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	ginkrakend "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
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
			"roles": "admin",
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

	signerEndpointCfg := newSignerEndpointCfg("HS256", "sim2", "http://localhost:5555/jwk/symmetric.json")

	buf := new(bytes.Buffer)
	logger, _ := logging.NewLogger("DEBUG", buf, "")
	hf := ginkrakendjose.HandlerFactory(ginkrakend.EndpointHandler, logger, nil)

	engine := gin.New()

	engine.POST(signerEndpointCfg.Endpoint, hf(signerEndpointCfg, tokenIssuer))

	fmt.Println("token request")
	req := httptest.NewRequest("POST", signerEndpointCfg.Endpoint, new(bytes.Buffer))

	ws := httptest.NewRecorder()
	engine.ServeHTTP(ws, req)

	responseData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Expiration   int    `json:"exp"`
	}{}
	json.Unmarshal(ws.Body.Bytes(), &responseData)

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

func tokenIssuer(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
	return &proxy.Response{
		Data: map[string]interface{}{
			"access_token": map[string]interface{}{
				"aud":   "http://api.example.com",
				"iss":   "http://example.com",
				"sub":   "1234567890qwertyuio",
				"jti":   "mnb23vcsrt756yuiomnbvcx98ertyuiop",
				"roles": []string{"user", "admin"},
				"exp":   1735689600,
			},
			"refresh_token": map[string]interface{}{
				"aud": "http://api.example.com",
				"iss": "http://example.com",
				"sub": "1234567890qwertyuio",
				"jti": "mnb23vcsrt756yuiomn12876bvcx98ertyuiop",
				"exp": 1735689600,
			},
			"exp": 1735689600,
		},
		Metadata: proxy.Metadata{
			StatusCode: 201,
		},
		IsComplete: true,
	}, nil
}

func newSignerEndpointCfg(alg, ID, URL string) *config.EndpointConfig {
	return &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/refresh-token",
		Method:   "POST",
		Backend: []*config.Backend{
			{
				URLPattern: "/refresh-token",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			krakendjose.SignerNamespace: map[string]interface{}{
				"alg":                  alg,
				"kid":                  ID,
				"jwk-url":              URL,
				"keys-to-sign":         []string{"access_token", "refresh_token"},
				"disable_jwk_security": true,
				"cache":                true,
			},
		},
	}
}
