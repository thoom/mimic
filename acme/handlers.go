package acme

import (
	"io"
	"net/http"
	"strconv"

	"github.com/thoom/mimic/mimic"
	"github.com/valyala/fasttemplate"
)

func DiscoveryHandler(w http.ResponseWriter, r *http.Request, cfg *mimic.Config) {
	template := `{
	"XceUS6czIlw": "https://community.letsencrypt.org/t/adding-random-entries-to-the-directory/33417",
	"keyChange": "http://{{HOST}}:{{PORT}}/acme/key-change",
	"meta": {
		"caaIdentities": [
		"happy-hacker-ca.invalid"
		],
		"termsOfService": "https://{{HOST}}:4431/terms/v7",
		"website": "https://github.com/letsencrypt/{{HOST}}"
	},
	"newAccount": "http://{{HOST}}:{{PORT}}/acme/new-acct",
	"newNonce": "http://{{HOST}}:{{PORT}}/acme/new-nonce",
	"newOrder": "http://{{HOST}}:{{PORT}}/acme/new-order",
	"revokeCert": "http://{{HOST}}:{{PORT}}/acme/revoke-cert"
}`

	t := fasttemplate.New(template, "{{", "}}")
	json := t.ExecuteString(map[string]interface{}{
		"HOST": cfg.Host,
		"PORT": strconv.Itoa(cfg.Port),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, json)
}
