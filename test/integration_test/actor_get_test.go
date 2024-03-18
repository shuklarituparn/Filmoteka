package integration_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/charmbracelet/log"
//	"io"
//	"net/http"
//	"testing"
//
//	"github.com/ozontech/allure-go/pkg/framework/provider"
//	"github.com/ozontech/allure-go/pkg/framework/suite"
//)
//
//type authToken struct {
//	Token string `json:"access_token"`
//}
//
//var token authToken
//
//type GettingActorsSuite struct {
//	suite.Suite
//}
//
//func (s *GettingActorsSuite) TestTokenGen(t provider.T) {
//	t.Epic("Getting Actors")
//	t.Feature("Login")
//	t.Title("Login in the system")
//	t.Description(`
//	In this step we will make the login request and get the JWT token`)
//
//	t.Tags("Login", "Get")
//
//	t.WithNewStep("Logging in", func(ctx provider.StepCtx) {
//		var requestbody = struct {
//			Email    string `json:"email"`
//			Password string `json:"password"`
//		}{
//			Email:    "admin@example.com",
//			Password: "adminpassword",
//		}
//		var buf bytes.Buffer
//		err := json.NewEncoder(&buf).Encode(requestbody)
//		if err != nil {
//			log.Fatal(err)
//		}
//		req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/users/login", &buf)
//		resp, err := http.DefaultClient.Do(req)
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer func(Body io.ReadCloser) {
//			err := Body.Close()
//			if err != nil {
//				t.Error("the response body could not be closed")
//			}
//		}(resp.Body)
//
//		body, _ := io.ReadAll(resp.Body)
//		if err := json.Unmarshal(body, &token); err != nil {
//			log.Error(err)
//		}
//
//		fmt.Println(token.Token)
//		t.WithNewStep("Get the list of actors", func(Ctx provider.StepCtx) {
//			t.Feature("Getting the list of actors")
//			authHeader := fmt.Sprintf("Bearer %s", token.Token)
//			getURL := fmt.Sprintf("http://localhost:8080/api/v1/actors/all?page=%d&page_size=%d", 1, 2)
//			newReq, _ := http.NewRequest("GET", getURL, nil)
//			newReq.Header.Set("Authorization", authHeader)
//			resp, err := http.DefaultClient.Do(newReq)
//			if err != nil {
//				log.Fatal(err)
//			}
//			defer func(Body io.ReadCloser) {
//				err := Body.Close()
//				if err != nil {
//					t.Error("the response body could not be closed")
//
//				}
//			}(resp.Body)
//
//			body, _ := io.ReadAll(resp.Body)
//			fmt.Println(string(body))
//
//		})
//	})
//}
//
//func TestStepTree(t *testing.T) {
//	t.Parallel()
//	suite.RunSuite(t, new(GettingActorsSuite))
//}
