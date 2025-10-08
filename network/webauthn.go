package network

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

func NewWebAuthn() *webauthn.WebAuthn {
	wconfig := &webauthn.Config{
		RPDisplayName: "APTrust",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	}
	webauthn, err := webauthn.New(wconfig)
	if err != nil {
		return nil
	}
	return webauthn
}
