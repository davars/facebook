package facebook

import (
	"encoding/json"
	"fmt"
	"io"
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

type Response interface {
	Error () error
}

type Map map[string]interface {}
func (resp Map) Error () (err error) {
	if resp["error"] == nil {
		return nil
	}

	fbError := resp["error"].(Map)
	return fmt.Errorf("Facebook Error (%v): %v", fbError["code"], fbError["message"])
}

// Ensures that we got a JSON map from Facebook
func checkMap(resp Response, inErr error) (m Map, outErr error) {
	if inErr != nil {
		return nil, inErr
	}
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	m, ok := resp.(Map)
	if !ok {
		return m, fmt.Errorf("Expected a JSON map, got %q instead", resp)
	}
	return
}

type Bool bool
func (resp Bool) Error () (err error) {
	return nil
}

// Ensures that we got a JSON bool from Facebook
func checkBool(resp Response, inErr error) (m Bool, outErr error) {
	if inErr != nil {
		return false, inErr
	}
	if resp.Error() != nil {
		return false, resp.Error()
	}
	m, ok := resp.(Bool)
	if !ok {
		return false, fmt.Errorf("Expected a JSON bool, got %q instead", resp)
	}
	return
}

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

	var decoded interface {}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&decoded)
	if err != nil {
		return nil, fmt.Errorf("JSON Decode failed: %q", err)
	}

	if _, ok := decoded.(map[string]interface {}); ok {
		return Map(decoded.(map[string]interface {})), nil
	}

	if _, ok := decoded.(bool); ok {
		return Bool(decoded.(bool)), nil
	}

	return
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
