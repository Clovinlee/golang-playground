package controllers

import (
	"bytes"
	"chris/gochris/initializers"
	"chris/gochris/models"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	const baseUrl string = "http://localhost:3000"
	const registerUrl string = "/register"
	const loginUrl string = "/login"

	var user models.User
	// var user models.User = models.User{Email: "test99@gmail.com", Password: "testpassword", Name: "Test User"}

	t.Run("Register", func(t *testing.T) {

		if user.Email != "" {
			t.Skip("Test skipped.. (User already registered)")
			return
		}

		t.Run("Fails", func(t *testing.T) {
			t.Run("Null Body", func(t *testing.T) {
				req, _ := http.NewRequest("POST", baseUrl+registerUrl, nil)
				req.Header.Set("Content-Type", "application/json; charset=utf-8")

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code: %d, got %d", http.StatusBadRequest, res.StatusCode)
				}

				defer res.Body.Close()

			})

			t.Run("Empty Email", func(t *testing.T) {
				payload := map[string]string{"password": "testpassword"}
				jsonPayload, _ := json.Marshal(payload)

				req, _ := http.NewRequest("POST", baseUrl+registerUrl, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code: %d, got %d", http.StatusBadRequest, res.StatusCode)
				}

				defer res.Body.Close()

			})

			t.Run("Empty Password", func(t *testing.T) {
				payload := map[string]string{"email": "test@g.com"}
				jsonPayload, _ := json.Marshal(payload)

				req, _ := http.NewRequest("POST", baseUrl+registerUrl, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code: %d, got %d", http.StatusBadRequest, res.StatusCode)
				}
				defer res.Body.Close()

			})

		})

		t.Run("Success", func(t *testing.T) {
			payload := map[string]string{"email": "test99@gmail.com", "password": "testpassword", "name": "Test User"}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", baseUrl+registerUrl, bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")

			client := &http.Client{}
			res, _ := client.Do(req)

			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected status code: %d, got %d", http.StatusBadRequest, res.StatusCode)
			} else {
				user.Email = payload["email"]
				user.Password = payload["password"] // Unhashed
				user.Name = payload["name"]
			}

			defer res.Body.Close()

		})

		// Email exist on register check
		defer func(t *testing.T) {
			t.Run("Register On Existing Email", func(t *testing.T) {
				if user.Email == "" {
					t.Skip("Test user not registered, test skipped..")
					return
				}

				payload := map[string]string{"email": user.Email, "password": user.Password}
				jsonPayload, _ := json.Marshal(payload)

				req, _ := http.NewRequest("POST", baseUrl+registerUrl, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")

				client := &http.Client{}
				res, _ := client.Do(req)

				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code: %d, got %d", http.StatusBadRequest, res.StatusCode)
				}

				defer res.Body.Close()
			})
		}(t)
	})

	t.Run("Login", func(t *testing.T) {
		expectedStatus := http.StatusOK

		payload := map[string]string{"email": user.Email, "password": user.Password}
		jsonPayload, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", baseUrl+loginUrl, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		client := &http.Client{}
		res, _ := client.Do(req)

		if res.StatusCode != expectedStatus {
			t.Errorf("Expected status code: %d, got %d", expectedStatus, res.StatusCode)
		}
	})

	t.Cleanup(func() {
		initializers.LoadEnvVariables("../.env")
		initializers.ConnectToDB()

		if user.Email == "" {
			return
		}

		// remove user registered
		result := initializers.DB.Delete(&models.User{}, "email = ?", user.Email)

		if result.RowsAffected > 0 {
			t.Log("User test cleanup successfully (User testing deleted)")
		} else {
			t.Error("Failed to cleanup user test (User testing not deleted)")
		}
	})
}
