package integration_test

import (
	"context"
	"errors"
	"github.com/ozontech/cute"
	"net/http"
	"testing"
	"time"
)

var Loginrequestbody = `{
		"email": "admin@example.com",
		"password": "adminpassword"
}`

func TestUserLogin(t *testing.T) {
	cute.NewTestBuilder().
		Title("Checking user login").
		Tags("user_login", "token", "auth").
		Feature("Login").
		Description("Testing if user can login into the system").
		CreateStep("Login").
		RequestBuilder(
			cute.WithMethod(http.MethodPost),
			cute.WithURI("https://api.rtprnshukla.ru/api/v1/users/login"),
			cute.WithHeaders(map[string][]string{
				"Content-Type": []string{"application/json"},
			}),
			cute.WithBody([]byte(Loginrequestbody)),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertResponse(
			func(resp *http.Response) error {
				if resp.ContentLength == 0 {
					return errors.New("content length is zero")
				}
				return nil
			},
		).ExecuteTest(context.Background(), t)
}
