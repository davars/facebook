package facebook

import (
	"net/url"
)

func (app Application) TestUsers() (users []Map, err error) {
	resp, err := checkMap(app.AppToken.Get("/"+app.Id+"/accounts/test-users", nil))
	if err != nil {
		return
	}

	for _, user := range resp["data"].([]interface{}) {
		users = append(users, Map(user.(map[string]interface{})))
	}

	return
}

func (app Application) CreateTestUser(values url.Values) (user Map, err error) {
	return checkMap(app.AppToken.Post("/"+app.Id+"/accounts/test-users", values, nil))
}

func (app Application) DeleteTestUser(id string) (success Bool, err error) {
	values := make(url.Values)
	values.Add("method", "delete")
	return checkBool(app.AppToken.Get("/"+id, values))
}
