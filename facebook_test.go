package facebook

import (
	"encoding/json"
	"io/ioutil"
    "os"
    "path"
    "testing"
)

var testApp *Application
func getTestApp(t *testing.T) (app *Application) {
	if testApp != nil {
		return testApp
	}

    filename := path.Join(path.Join(os.Getenv("PWD"), "tests"), "fb_app.json")
    config, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Unable to open config file You need to create a ./tests/fb_app.json file with your app's credentials to run the tests.")
	}

    app = new(Application)
    err = json.Unmarshal(config, &app)
	if err != nil {
		t.Fatalf("Unable to parse config file: %q", err)
	}

	testApp = app
    return testApp
}

var testUser Response
func getTestUser(t *testing.T) (user Response) {
	if testUser != nil {
		return testUser
	}
	app := getTestApp(t)

	users, err := app.TestUsers()
	if err != nil {
		t.Fatalf("Unable to get test users: %q (Response: %v)", err, users)
	}

	if len(users) == 0 {
		t.Fatalf("Please create a test user to run these tests.")
	}

	testUser = users[0]
	return testUser
}

func TestGet(t *testing.T) {
	app := getTestApp(t)
	resp, err := Get("/" + string(app.Id), nil)
	if(err != nil) {
		t.Fatalf("Get failed: %q", err)
	}
	if resp["id"] != app.Id {
		t.Errorf("Got the wrong id from the Graph Api (this should never happen).  Expected %q,  got %q. (Response: %q)", app.Id, resp["id"], resp)
	}
}

func TestCreateTestUser(t *testing.T) {
	user := getTestUser(t)
	t.Logf("%v", user)
}

