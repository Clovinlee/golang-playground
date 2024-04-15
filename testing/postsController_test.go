package controllers

import (
	"bytes"
	"chris/gochris/models"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
)

const baseUrl string = "http://localhost:3000"
const postUrl string = "/posts"
const loginUrl string = "/login"

var authCookie *http.Cookie

func TestPost(t *testing.T) {
	// var authCookie *http.Cookie

	doLogin(t)

	if authCookie == nil {
		t.Fatal("Authcookie undefined (nil)")
	}
	var postCreated models.Post

	t.Run("Create", func(t *testing.T) {
		t.Run("Fails", func(t *testing.T) {

			t.Run("Unauthorized Access", func(t *testing.T) {
				expectedStatus := http.StatusUnauthorized

				req, _ := http.NewRequest("POST", baseUrl+postUrl, nil)
				req.Header.Set("Content-Type", "application/json; charset=utf-8")

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != expectedStatus {
					t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
				}
				defer res.Body.Close()
			})

			t.Run("Empty Title", func(t *testing.T) {
				expectedStatus := http.StatusBadRequest

				payload := map[string]string{"body": "This is body of test post"}
				jsonPayload, _ := json.Marshal(payload)

				req, _ := http.NewRequest("POST", baseUrl+postUrl, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				req.AddCookie(authCookie)

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != expectedStatus {
					t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
				}

				defer res.Body.Close()
			})

			t.Run("Empty Body", func(t *testing.T) {
				expectedStatus := http.StatusBadRequest

				payload := map[string]string{"title": "This is title of test post"}
				jsonPayload, _ := json.Marshal(payload)

				req, _ := http.NewRequest("POST", baseUrl+postUrl, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				req.AddCookie(authCookie)

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != expectedStatus {
					t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
				}

				defer res.Body.Close()

			})
		})

		t.Run("Success", func(t *testing.T) {
			expectedStatus := http.StatusOK

			payload := map[string]string{"title": "This is title of test post SUCCESS", "body": "This is body of test post SUCCESS"}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", baseUrl+postUrl, bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.AddCookie(authCookie)

			client := &http.Client{}
			res, _ := client.Do(req)

			if res.StatusCode != expectedStatus {
				t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
			}

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error("Error reading response body: ", err)
			}

			type postCreateResponse struct {
				Post models.Post
			}

			var result postCreateResponse
			json.Unmarshal([]byte(resBody), &result)

			postCreated = result.Post
		})
	})

	t.Run("Delete", func(t *testing.T) {
		expectedStatus := http.StatusOK

		req, _ := http.NewRequest("DELETE", baseUrl+postUrl+"/"+strconv.Itoa(int(postCreated.ID)), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.AddCookie(authCookie)

		client := &http.Client{}
		res, _ := client.Do(req)

		if res.StatusCode != expectedStatus {
			t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
		}
		defer res.Body.Close()
	})

}

func doLogin(t *testing.T) {
	// test account
	var user models.User = models.User{Email: "testposts@gmail.com", Password: "testpassword"}

	payload := map[string]string{"email": user.Email, "password": user.Password}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", baseUrl+loginUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, _ := client.Do(req)

	if res.StatusCode != http.StatusOK {
		t.Fatal("Login failed")
	} else {
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "Authorization" {
				authCookie = cookie
			}
		}
		defer res.Body.Close()
	}
}
