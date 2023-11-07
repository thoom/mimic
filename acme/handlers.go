package acme

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/thoom/mimic/mimic"
)

type DirectoryJsonMeta struct {
	CaaIdentities  []string `json:"caaIdentities"`
	TermsOfService string   `json:"termsOfService"`
	Website        string   `json:"website"`
}

type DirectoryJson struct {
	KeyChange  string            `json:"keyChange"`
	Meta       DirectoryJsonMeta `json:"meta"`
	NewAccount string            `json:"newAccount"`
	NewNonce   string            `json:"newNonce"`
	NewOrder   string            `json:"newOrder"`
	RevokeCert string            `json:"revokeCert"`
}

func DirectoryHandler(w http.ResponseWriter, r *http.Request, cfg *mimic.Config) {
	dj := DirectoryJson{
		KeyChange: cfg.HostURL + "/acme/key-change",
		Meta: DirectoryJsonMeta{
			CaaIdentities:  []string{"mimic-ca.invalid"},
			TermsOfService: cfg.HostURL + "/terms/v1",
			Website:        "https://github.com/thoom/mimic",
		},
		NewAccount: cfg.HostURL + "/acme/new-acct",
		NewNonce:   cfg.HostURL + "/acme/new-nonce",
		NewOrder:   cfg.HostURL + "/acme/new-order",
		RevokeCert: cfg.HostURL + "/acme/revoke-cert",
	}

	response, _ := json.Marshal(dj)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(response))
}
