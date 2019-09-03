package main

import (
	"bytes"
	"context"
	"time"

	krakendjose "github.com/devopsfaith/krakend-jose"
	ginkrakendjose "github.com/devopsfaith/krakend-jose/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	ginkrakend "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	signerEndpointCfg := newSignerEndpointCfg("HS256", "sim2", "http://localhost:5555/jwk/symmetric.json")

	buf := new(bytes.Buffer)
	logger, _ := logging.NewLogger("DEBUG", buf, "")
	hf := ginkrakendjose.HandlerFactory(ginkrakend.EndpointHandler, logger, nil)

	engine := gin.New()
	engine.POST(signerEndpointCfg.Endpoint, hf(signerEndpointCfg, tokenIssuer))

	engine.Run(":2222")
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
