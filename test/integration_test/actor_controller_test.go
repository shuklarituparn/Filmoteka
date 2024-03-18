package integration_test

//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/ozontech/cute"
//	"io"
//	"net/http"
//	"testing"
//	"time"
//)
//
//var requestbody = `{
//		"email": "admin@example.com",
//		"password": "adminpassword"
//}`
//
//type Auth struct {
//	Token string `json:"access_token"`
//}
//
//var token Auth
//var authTokenString string
//
//func TestExampleTest(t *testing.T) {
//	getURL := fmt.Sprintf("http://localhost:8080/api/v1/actors/all?page=%d&page_size=%d", 1, 2)
//	cute.NewTestBuilder().
//		Title("Getting the actors").
//		Tags("one_step", "some_local_tag", "json").
//		Feature("some_feature").
//		Description("some_description").
//		CreateStep("SFFS").
//		RequestBuilder(
//			cute.WithMethod(http.MethodPost),
//			cute.WithURI("http://localhost:8080/api/v1/users/login"),
//			cute.WithHeaders(map[string][]string{
//				"Content-Type": []string{"application/json"},
//			}),
//			cute.WithBody([]byte(requestbody)),
//		).
//		ExpectExecuteTimeout(10*time.Second).
//		ExpectStatus(http.StatusOK).
//		AssertResponse(
//			func(resp *http.Response) error {
//				if resp.ContentLength == 0 {
//					return errors.New("content length is zero")
//				}
//				b, err := io.ReadAll(resp.Body)
//				if err != nil {
//					return err
//				}
//				if err := json.Unmarshal(b, &token); err != nil {
//					fmt.Print(err)
//				}
//				fmt.Println(token.Token)
//				authTokenString = fmt.Sprintf("Bearer %s", token.Token)
//				return nil
//
//			},
//		).
//		NextTest().
//		Create().
//		RequestBuilder(
//			cute.WithMethod(http.MethodGet),
//			cute.WithURI(getURL),
//			cute.WithHeaders(map[string][]string{
//				"Authorization": []string{authTokenString},
//			}),
//		).
//		ExpectStatus(http.StatusOK).
//		ExecuteTest(context.Background(), t)
//
//	fmt.Println(authTokenString)
//
//}
