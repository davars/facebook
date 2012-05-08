package facebook

import (
	"net/url"
    "testing"
)

func TestAccessToken(t *testing.T) {
	app := getTestApp(t)
	resp, err := app.AccessToken(url.Values {"grant_type": []string{"client_credentials"}})
	checkFatal(t, err)

	values, err := url.ParseQuery(resp)
	checkFatal(t, err)
	
	if AccessToken(values["access_token"][0]) != app.AppToken {
		t.Fatalf("access_tokens don't match")
	}
}

