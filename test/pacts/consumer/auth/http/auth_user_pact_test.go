package auth_user_pact_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

// DTOs for the user-service verify endpoint

type verifyUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userDTO struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// Very small HTTP client that auth-service would use to call user-service.
type UserVerificationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserVerificationClient(baseURL string) *UserVerificationClient {
	return &UserVerificationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// VerifyCredentials calls POST /internal/users/verify on user-service.
func (c *UserVerificationClient) VerifyCredentials(email, password string) (*userDTO, error) {
	reqBody := verifyUserRequest{
		Email:    email,
		Password: password,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+"/internal/users/verify",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("close response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var u userDTO
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("decode body: %w", err)
	}

	return &u, nil
}

func TestVerifyUserCredentialsPact(t *testing.T) {
	t.Parallel()

	// Create the Pact mock provider (user-service) using the **V4** mock provider
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "auth-service",
		Provider: "user-service",

		// Adjust to your repo layout as needed
		PactDir: "./pacts",
		LogDir:  "./logs",
	})
	assert.NoError(t, err)

	email := "user@example.com"
	password := "super-secret-password"

	// Arrange: define the interaction (contract) between auth-service and user-service
	mockProvider.
		AddInteraction().
		Given("a user exists with this email and password").
		UponReceiving("a request to verify valid user credentials").
		WithRequest(http.MethodPost, "/internal/users/verify",
			func(b *consumer.V4RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Accept", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"email":    matchers.Like(email),
					"password": matchers.Like(password),
				})
			},
		).
		WillRespondWith(200,
			func(r *consumer.V4ResponseBuilder) {
				r.Header("Content-Type", matchers.Regex(
					"application/json",
					`^application/json($|;.*)$`,
				))
				r.JSONBody(map[string]interface{}{
					"id":     matchers.Like("8a26b19d-8a33-4ece-87b1-7b7c2fb9e0ad"),
					"email":  matchers.Like(email),
					"status": matchers.Like("active"),
				})
			},
		)

	// Act + Assert: run the client against the Pact mock server
	err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
		client := NewUserVerificationClient(baseURL)

		user, err := client.VerifyCredentials(email, password)
		assert.NoError(t, err)
		if err != nil {
			// If the client returned an error, just propagate it so Pact fails quickly
			return err
		}

		assert.Equal(t, email, user.Email)
		assert.Equal(t, "active", user.Status)
		assert.NotEmpty(t, user.ID)

		return nil // important: returning nil tells Pact the interaction succeeded
	})
	assert.NoError(t, err)
}
