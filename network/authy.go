package network

import (
	"errors"
	stdlog "log"
	"net/url"
	"os"
	"time"

	"github.com/dcu/go-authy"
	"github.com/rs/zerolog"
)

var authyLoginMessage = "Log in to the APTrust registry."
var authyTimeout = (45 * time.Second)

type AuthyClient struct {
	client  *authy.Authy
	enabled bool
	log     zerolog.Logger
}

func NewAuthyClient(authyEnabled bool, authyAPIKey string, log zerolog.Logger) *AuthyClient {
	// Authy library logs to Stderr by default. We want either our logger
	// (which doesn't support the standard go logger interface) or Stdout
	// because Docker gathers Stdout logs.
	authy.Logger = stdlog.New(os.Stdout, "[authy] ", stdlog.LstdFlags)
	return &AuthyClient{
		client:  authy.NewAuthyAPI(authyAPIKey),
		enabled: authyEnabled,
		log:     log,
	}
}

// AwaitOneTouch sends a OneTouch login request via Authy and awaits
// the user's response. Param authyID is the user's AuthyID. Param
// userEmail is used for logging.
//
// If request is approved, this returns true. Otherwise, false.
func (ac *AuthyClient) AwaitOneTouch(userEmail, authyID string) (bool, error) {
	if !ac.enabled {
		return false, errors.New("Cannot send OneTouch because Authy client is disabled.")
	}
	details := authy.Details{}
	req, err := ac.client.SendApprovalRequest(authyID, authyLoginMessage, details, url.Values{})
	if err != nil {
		ac.log.Error().Msgf("AuthyOneTouch error for %s: %v", userEmail, err)
		return false, err
	}
	ac.log.Info().Msgf("AuthyOneTouch request id for %s: %s", userEmail, req.UUID)
	status, err := ac.client.WaitForApprovalRequest(req.UUID, authyTimeout, url.Values{})
	if status == authy.OneTouchStatusApproved {
		ac.log.Info().Msgf("AuthyOneTouch approved for %s (%s)", userEmail, req.UUID)
		return true, nil
	} else {
		ac.log.Warn().Msgf("AuthyOneTouch %s for %s (%s)", status, userEmail, req.UUID)
	}
	return false, nil
}
