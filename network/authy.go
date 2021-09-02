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

// ErrAuthyDisabled means authy isn't enabled here. You can change that
// in the .env file.
var ErrAuthyDisabled = errors.New("authy is not enabled in this environment")

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
// This is a blocking request that waits up to 45 seconds for a user
// to approve the one-touch push notification.
//
// If request is approved, this returns true. Otherwise, false.
func (ac *AuthyClient) AwaitOneTouch(userEmail, authyID string) (bool, error) {
	if !ac.enabled {
		return false, ErrAuthyDisabled
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

// RegisterUser registers a user with Authy for this app. Note that
// users need separate registrations for each environment (dev, demo,
// prod, etc.).
//
// On success, this returns the user's new AuthyID. The caller is
// responsible for attaching that ID to the user object and saving
// it to the database.
//
// Use user.CountryCodeAndPhone() to get country code and phone number,
// as these need to be separate. Do not pass user.PhoneNumber in format
// "+<country_code><number>" because that won't work.
func (ac *AuthyClient) RegisterUser(userEmail string, countryCode int, phone string) (string, error) {
	authyUser, err := ac.client.RegisterUser(userEmail, countryCode, phone, url.Values{})
	if err != nil {
		ac.log.Error().Msgf("Can't register user %s (%d %s) with Authy: %v", userEmail, countryCode, phone, err)
		return "", err
	}
	return authyUser.ID, err
}
