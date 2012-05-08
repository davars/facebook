package facebook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
)

func (app Application) ParseSignedRequest(r string) (parsed Map, err error) {
	parts := strings.Split(r, ".")

	payload, err := base64.URLEncoding.DecodeString(pad64(parts[1]))
	if err != nil {
		return nil, errors.New("Malformed Signed Request")
	}

	err = json.Unmarshal(payload, &parsed)
	if err != nil {
		return nil, errors.New("Malformed Signed Request")
	}

	if parsed["algorithm"] != "HMAC-SHA256" {
		return nil, errors.New("Unknown algorithm " + parsed["algorithm"].(string))
	}

	hmac := hmac.New(sha256.New, []byte(app.Secret))
	hmac.Write([]byte(parts[1]))

	sig := base64.URLEncoding.EncodeToString(hmac.Sum(nil))

	if pad64(parts[0]) != sig {
		return nil, errors.New("Bad Signature")
	}

	return
}

func (app Application) AccessToken(values url.Values) (response string, err error) {
	if values == nil {
		return "", errors.New("Missing extra form paramters, perhaps redirect_uri or grant_type?")
	}
	values.Add("client_id", app.Id)
	values.Add("client_secret", app.Secret)
	r, err := graphRequest("GET", "/oauth/access_token", values, nil)
	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		return "", err
	}
	
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
