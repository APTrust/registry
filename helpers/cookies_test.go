package helpers_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cookieUser *pgmodels.User

type TestCookieSetter struct {
	Cookies map[string]http.Cookie
}

func NewCookieSetter() *TestCookieSetter {
	return &TestCookieSetter{
		Cookies: make(map[string]http.Cookie),
	}
}

// SetCookie adds a cookie to the Cookies collection.
func (c *TestCookieSetter) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.Cookies[name] = http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteStrictMode,
	}
}

// SetSameSite is an empty mock
func (c *TestCookieSetter) SetSameSite(samesite http.SameSite) {
	// No op
}

// Get is a mock function to return user for cookie tests.
func (c *TestCookieSetter) Get(key string) (value interface{}, exists bool) {
	return cookieUser, true
}

func initCookieUser(t *testing.T) {
	if cookieUser == nil {
		var err error
		cookieUser, err = pgmodels.UserByEmail("admin@inst1.edu")
		require.Nil(t, err)
		require.NotNil(t, cookieUser)
	}
}

func getSetter(t *testing.T) *TestCookieSetter {
	initCookieUser(t)
	return NewCookieSetter()
}

func TestSessionCookie(t *testing.T) {
	setter := getSetter(t)
	helpers.SetSessionCookie(setter, cookieUser)

	name := common.Context().Config.Cookies.SessionCookie
	cookie := setter.Cookies[name]
	testCookieWasSet(t, cookie, name)

	helpers.DeleteSessionCookie(setter)
	cookie = setter.Cookies[name]
	testCookieWasDeleted(t, cookie, name)
}

func TestFlashCookie(t *testing.T) {
	setter := getSetter(t)
	helpers.SetFlashCookie(setter, "Flash!")

	name := common.Context().Config.Cookies.FlashCookie
	cookie := setter.Cookies[name]
	testCookieWasSet(t, cookie, name)

	helpers.DeleteFlashCookie(setter)
	cookie = setter.Cookies[name]
	testCookieWasDeleted(t, cookie, name)
}

func TestPrefsCookie(t *testing.T) {
	setter := getSetter(t)
	helpers.SetPrefsCookie(setter, "Prefs!")

	name := common.Context().Config.Cookies.PrefsCookie
	cookie := setter.Cookies[name]
	testCookieWasSet(t, cookie, name)

	helpers.DeletePrefsCookie(setter)
	cookie = setter.Cookies[name]
	testCookieWasDeleted(t, cookie, name)
}

func TestCSRFCookie(t *testing.T) {
	setter := getSetter(t)
	token, err := helpers.SetCSRFCookie(setter)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)

	name := constants.CSRFCookieName
	cookie := setter.Cookies[name]
	testCookieWasSet(t, cookie, name)

	helpers.DeleteCSRFCookie(setter)
	cookie = setter.Cookies[name]
	testCookieWasDeleted(t, cookie, name)
}

func TestCurrentUser(t *testing.T) {
	setter := getSetter(t)
	user := helpers.CurrentUser(setter)
	require.NotNil(t, user)
	assert.Equal(t, cookieUser.ID, user.ID)
	assert.Equal(t, cookieUser.Email, user.Email)
}

func testCookieWasSet(t *testing.T, cookie http.Cookie, name string) {
	require.NotNil(t, cookie, name)
	config := common.Context().Config.Cookies
	assert.Equal(t, name, cookie.Name)
	assert.NotEmpty(t, cookie.Value)
	assert.Equal(t, config.MaxAge, cookie.MaxAge)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, config.Domain, cookie.Domain)
	assert.Equal(t, config.HTTPSOnly, cookie.Secure)
	assert.True(t, cookie.HttpOnly)
}

func testCookieWasDeleted(t *testing.T, cookie http.Cookie, name string) {
	require.NotNil(t, cookie, name)
	assert.Equal(t, -1, cookie.MaxAge)
}
