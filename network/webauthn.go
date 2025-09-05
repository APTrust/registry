package network

import (
	"errors"
	stdlog "log"
	"net/url"
	"os"
	"time"
)

func NewWebAuthn() WebAuthn {
	wconfig := &webAuthn.Config{
		RPDisplayName: "APTrust"
		RPID: "localhost",
		RPOrigins: []string{"http://localhost:8080"}
	}
	webauthn, err := webauthn.New(wconfig)
	if err != nil {
		return nil
	}
	return webauthn
}
