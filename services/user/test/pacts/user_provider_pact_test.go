package user_provider_pact_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	servergen "github.com/incheat/go-playground/services/user/internal/api/gen/oapi/private/server"
	userhandler "github.com/incheat/go-playground/services/user/internal/handler/http"
	userservice "github.com/incheat/go-playground/services/user/internal/service/user"
	"github.com/incheat/go-playground/services/user/pkg/model"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------
// Fake repository + service wiring (adjust to your real interfaces)
// -------------------------------------------------------------------

type fakeUserRepo struct {
	email    string
	password string
}

func (f *fakeUserRepo) GetUserByEmail(_ context.Context, email string) (*model.User, error) {
	user := &model.User{
		ID:           "8a26b19d-8a33-4ece-87b1-7b7c2fb9e0ad",
		Email:        f.email,
		Status:       "active",
		PasswordHash: f.password,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}
	if email == user.Email {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (f *fakeUserRepo) CreateUser(_ context.Context, _ string, _ *model.User) error {
	return nil
}

// -------------------------------------------------------------------
// Provider Pact Test (matches consumer pact with 200 + 401 interactions)
// -------------------------------------------------------------------

func TestUserServiceProviderPact(t *testing.T) {

	root, filepathErr := repoRoot()
	require.NoError(t, filepathErr)
	pactFile := filepath.Join(root, "test/pacts/consumer/auth/http/pacts/auth-service-user-service.json")

	// Fake repo we can seed via provider states
	repo := &fakeUserRepo{}

	// Real service + real HTTP handlers/router (adjust ctor signatures if needed)
	service := userservice.New(repo)
	userImpl := userhandler.New(service)

	// Real HTTP server, ephemeral port, no goroutine management
	strict := servergen.NewStrictHandler(userImpl, nil)
	apiHandler := servergen.HandlerFromMux(strict, chi.NewRouter())
	server := httptest.NewServer(apiHandler)
	defer server.Close()

	verifier := provider.NewVerifier()

	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		Provider:        "user-service",
		ProviderBaseURL: server.URL,
		PactFiles: []string{
			pactFile,
		},

		// Must include ALL Given(...) states present in the pact file
		StateHandlers: models.StateHandlers{
			// Matches the 200 interaction's Given(...)
			"a user exists with this email and password": func(setup bool, _ models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					repo.email = "user@example.com"
					repo.password = "super-secret-password"
				}
				return nil, nil
			},

			// Matches the 401 interaction's Given(...)
			// We seed the same "real" credentials; the pact request will use a wrong password,
			// so Verify() will return ErrInvalidCredentials -> handler should map to 401.
			"invalid user credentials": func(setup bool, _ models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					repo.email = "user@example.com"
					repo.password = "wrong-password"
				}
				return nil, nil
			},
		},

		// Optional but often helpful: verify provider has the expected content type even on 401
		// If your consumer pact asserts Content-Type on 401, keep it consistent in handler.
		// Otherwise you can remove this section.
		EnablePending:              false,
		FailIfNoPactsFound:         true,
		PublishVerificationResults: false,
	})

	require.NoError(t, err)
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
