package facebook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Application struct {
	Id        string
	Secret    string
	Namespace string
	AppToken  AccessToken
}

type AccessToken string

type Response map[string]interface{}

func pad64(s string) string {
	padLength := 4 - len(s)%4
	if padLength == 4 {
		return s
	}

	return s + strings.Repeat("=", padLength)
}

const GRAPH_API = "https://graph.facebook.com"

func graphRequest(method string, id string, values url.Values, body io.Reader) (r *http.Response, err error) {
	url := GRAPH_API + id + "?" + values.Encode()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func jsonRequest(method string, id string, values url.Values, body io.Reader) (resp Response, err error) {
	r, err := graphRequest(method, id, values, body)
	if(r != nil) {
		defer r.Body.Close()	
	}
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(r.Body)
	resp = make(Response)
	err = dec.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return
}

func (resp Response) Error () (err error) {
	if resp["error"] == nil {
		return nil
	}

	fbError := resp["error"].(map[string]interface {})
	return fmt.Errorf("Facebook Error (%v): %v", fbError["code"], fbError["message"])
}

func Get(id string, values url.Values) (resp Response, err error) {
	return AccessToken("").Get(id, values)
}

func (token AccessToken) Get(id string, values url.Values) (resp Response, err error) {
	if values == nil {
		values = make(url.Values)
	}
	values.Add("access_token", string(token))

	return jsonRequest("GET", id, values, nil)
}


func (token AccessToken) Post(id string, values url.Values, body io.Reader) (resp Response, err error) {
	if values == nil {
		values = make(url.Values)
	}
	values.Add("access_token", string(token))

	return jsonRequest("POST", id, values, body)
}

func (app Application) ParseSignedRequest(r string) (parsed Response, err error) {
	parts := strings.Split(r, ".")

	payload, err := base64.URLEncoding.DecodeString(pad64(parts[1]))
	if err != nil {
		return nil, errors.New("Malformed Signed Request")
	}

	parsed = make(Response)
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

func (app Application) TestUsers() (users []Response, err error) {
	resp, err := app.AppToken.Get("/" + app.Id + "/accounts/test-users", nil)
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		return nil, resp.Error()
	}

	for _, user := range resp["data"].([]interface {}) {
		users = append(users, Response(user.(map[string]interface {})))
	}

	return 
}

func (app Application) CreateTestUser(values url.Values) (user Response, err error) {
	return app.AppToken.Post("/" + app.Id + "/accounts/test-users", values, nil)
}

func (app Application) AccessToken(values url.Values) (token AccessToken, err error) {
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

	return AccessToken(data), nil
}
