package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ozontech/cute"
	"io"
	"net/http"
	"testing"
	"time"
)

var requestbody = `{
		"email": "admin@example.com",
		"password": "adminpassword"
}`

type Auth struct {
	Token string `json:"access_token"`
}

var token Auth

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
			cute.WithBody([]byte(requestbody)),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertResponse(
			func(resp *http.Response) error {
				if resp.ContentLength == 0 {
					return errors.New("content length is zero")
				}
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(b, &token); err != nil {
					fmt.Print(err)
				}
				fmt.Println(token.Token)
				return nil
			},
		).ExecuteTest(context.Background(), t)
}
