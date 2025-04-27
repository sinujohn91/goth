package microsoftonline_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	provider := microsoftonlineProvider()

	a.Equal(provider.ClientKey, os.Getenv("MICROSOFTONLINE_KEY"))
	a.Equal(provider.Secret, os.Getenv("MICROSOFTONLINE_SECRET"))
	a.Equal(provider.CallbackURL, "/foo")
}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := microsoftonlineProvider()
	a.Implements((*goth.Provider)(nil), p)
}

func Test_BeginAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	provider := microsoftonlineProvider()
	session, err := provider.BeginAuth("test_state")
	s := session.(*microsoftonline.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "login.microsoftonline.com/common/oauth2/v2.0/authorize")
}

func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := microsoftonlineProvider()
	session, err := provider.UnmarshalSession(`{"AuthURL":"https://login.microsoftonline.com/common/oauth2/v2.0/authorize","AccessToken":"1234567890","ExpiresAt":"0001-01-01T00:00:00Z"}`)
	a.NoError(err)

	s := session.(*microsoftonline.Session)
	a.Equal(s.AuthURL, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	a.Equal(s.AccessToken, "1234567890")
}

func Test_BeginAuthWithPrompt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := microsoftonlineProvider()
	provider.SetPrompt("select_account")
	session, err := provider.BeginAuth("test_state")
	s := session.(*microsoftonline.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", os.Getenv("MICROSOFTONLINE_KEY")))
	a.Contains(s.AuthURL, "state=test_state")
	a.Contains(s.AuthURL, "prompt=select_account")
}

func Test_BeginAuthWithInvalidPrompt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := microsoftonlineProvider()
	provider.SetPrompt("invalid_prompt")
	session, err := provider.BeginAuth("test_state")
	s := session.(*microsoftonline.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", os.Getenv("MICROSOFTONLINE_KEY")))
	a.Contains(s.AuthURL, "state=test_state")
	a.NotContains(s.AuthURL, "prompt=select_account")
}

func Test_BeginAuthWithLoginHint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := microsoftonlineProvider()
	provider.SetLoginHint("john@example.com")
	session, err := provider.BeginAuth("test_state")
	s := session.(*microsoftonline.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", os.Getenv("MICROSOFTONLINE_KEY")))
	a.Contains(s.AuthURL, "login_hint=john%40example.com")
}

func microsoftonlineProvider() *microsoftonline.Provider {
	return microsoftonline.New(os.Getenv("MICROSOFTONLINE_KEY"), os.Getenv("MICROSOFTONLINE_SECRET"), "/foo")
}
