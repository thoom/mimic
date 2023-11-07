package mimic

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type JoseJson struct {
	Protected      string `json:"protected"`
	Signature      string `json:"signature"`
	EncodedPayload string `json:"payload"`
	Payload        string
	Jwt            struct {
		Alogrithm string `json:"alg"`
		Nonce     string `json:"nonce"`
		URL       string `json:"url"`
		Jwk       struct {
			E   string `json:"e"`
			Kty string `json:"kty"`
			N   string `json:"n"`
		}
	}
}

func (jose *JoseJson) DecodeProtected() {
	decoded, _ := base64.RawStdEncoding.DecodeString(jose.Protected)
	if err := json.Unmarshal([]byte(decoded), &jose.Jwt); err != nil {
		log.Fatal(err)
	}
}

func (jose *JoseJson) DecodePayload() {
	decoded, _ := base64.RawStdEncoding.DecodeString(jose.EncodedPayload)
	jose.Payload = string(decoded)
}
