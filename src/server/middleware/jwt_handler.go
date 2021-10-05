package middleware

import (
	"encoding/json"
	"errors"
	"insurance-ng/src/server/config"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

type jSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type jwks struct {
	Keys []jSONWebKeys `json:"keys"`
}

var auth0Jwks = jwks{}

// JwtMiddleware processes and verifis jwt before request is processed
func JwtMiddleware() *jwtmiddleware.JWTMiddleware {
	authAud := os.Getenv("AUTH0_CLIENT_ID")
	authDomain := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"

	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := authAud
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := authDomain
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token, authDomain)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		Debug:         config.IsDebugMode(),
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPemCert(token *jwt.Token, authDomain string) (string, error) {
	cert := ""
	// resp, err := http.Get(auth_domain + ".well-known/jwks.json")
	// defer resp.Body.Close()
	if auth0Jwks.Keys == nil {
		resp, err := os.Open(os.Getenv("AUTH0_JWT_WELKNOWNS_PATH"))
		if err != nil {
			return cert, err
		}
		defer resp.Close()

		err = json.NewDecoder(resp).Decode(&auth0Jwks)
		if err != nil {
			return cert, err
		}
	}

	for k := range auth0Jwks.Keys {
		if token.Header["kid"] == auth0Jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + auth0Jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
