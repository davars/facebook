package facebook

import (
	"net/url"
	"testing"
	"testing/quick"
)

func TestAccessToken(t *testing.T) {
	app := getTestApp(t)
	values, err := app.AccessToken(url.Values{"grant_type": []string{"client_credentials"}})
	checkFatal(t, err)

	if AccessToken(values["access_token"][0]) != app.AppToken {
		t.Fatalf("access_tokens don't match")
	}
}

func TestParseValidSignedRequest(t *testing.T) {
	example_sr := "vlXgu64BQGFSQrY0ZcJBZASMvYvTHu9GQ0YM9rjPSso.eyJhbGdvcml0aG0iOiJITUFDLVNIQTI1NiIsIjAiOiJwYXlsb2FkIn0"
	example_data := Map{"0": "payload"}
	example_app := Application{
		Secret: "secret",
	}
	parsed, err := example_app.ParseSignedRequest(example_sr)
	if err != nil {
		t.Fatalf("Signed request not parsed: %q", err)
	}
	for k, v := range example_data {
		if parsed[k].(string) != v {
			t.Errorf("Incorrect request content, expected %q, got %q", v, parsed[k])
		}
	}
}

func TestQuickCheckParseSignedRequest(t *testing.T) {
	example_app := Application{
		Secret: "secret",
	}

	f := func(sr string) bool {
		parsed, err := example_app.ParseSignedRequest(sr)
		if parsed != nil || err == nil {
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

}
