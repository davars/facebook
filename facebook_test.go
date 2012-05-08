package facebook

import (
	"encoding/json"
	"io/ioutil"
    "os"
    "path"
    "regexp"
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

	for _, user := range users {
		if Map(user)["access_token"] != nil {
			testUser = user
			break
		}
	}
	if testUser == nil {
		t.Fatalf("Please create a test user that has authorized your app to run these tests.")
	}

	return testUser
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Failed: %q", err)
	}
}

func TestGet(t *testing.T) {
	app := getTestApp(t)
	resp, err := Get("/" + string(app.Id), nil)
	if err != nil {
		t.Fatalf("Get failed: %q", err)
	}
	if resp.(Map)["id"] != app.Id {
		t.Errorf("Got the wrong id from the Graph Api (this should never happen).  Expected %q,  got %q. (Response: %q)", app.Id, resp.(Map)["id"], resp)
	}
}

func TestGetWithToken(t *testing.T) {
	user := getTestUser(t)
	resp, err := AccessToken(user.(Map)["access_token"].(string)).Get("/me", nil)
	if err != nil {
		t.Fatalf("Get failed: %q", err)
	}
	if resp.(Map)["id"] != user.(Map)["id"] {
		t.Errorf("Got the wrong id from the Graph Api (this should never happen).  Expected %q,  got %q. (Response: %q)", user.(Map)["id"], resp.(Map)["id"], resp)
	}
}

func TestCreateAndDeleteAndListTestUser(t *testing.T) {
	if !testing.Short() {
		app := getTestApp(t)
		user, err := app.CreateTestUser(nil)
		if err != nil {
			t.Fatalf("CreateTestUser failed: %q", err)
		}
		t.Logf("User created: %v", user)

		matches, err := regexp.MatchString("test_account_login", user["login_url"].(string))
		if err != nil {
			t.Fatalf("Regexp failed: %q", err)
		}
		if !matches {
			t.Errorf("Unexpected value for login_url %q", user["login_url"])
		}

		success, err := app.DeleteTestUser(user["id"].(string))
		if err != nil || !bool(success) {
			t.Errorf("Unable to delete test user: %q", err)
		}

		testUserList, err := app.TestUsers()

		for _, testUser := range testUserList {
			if user["id"] == testUser["id"] {
				t.Errorf("Test user not successfully deleted: %v", user)
			}
		}	
	}
}

