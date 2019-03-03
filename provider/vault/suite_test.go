package vault_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

func TestVault(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vault Suite")
}

func newAuthHandler() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("POST", "/v1/auth/kubernetes/login"),
		RespondWithJSONEncoded(http.StatusOK, &api.Secret{
			Auth: &api.SecretAuth{ClientToken: "my-token"},
		}),
	)
}

func newAuthHandlerFailed() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("POST", "/v1/auth/kubernetes/login"),
		RespondWithJSONEncoded(http.StatusInternalServerError,
			&api.ErrorResponse{
				Errors: []string{"oh no!"},
			},
		),
	)
}

func newGetMntHandler() http.HandlerFunc {
	type mount struct {
		Renewable bool            `json:"renewable"`
		Data      api.MountOutput `json:"data" mapstructure:"data"`
	}

	output := &mount{
		Data: api.MountOutput{
			Type: "kv",
			Options: map[string]string{
				"version": "2",
			},
		},
	}

	return CombineHandlers(
		VerifyRequest("GET", "/v1/sys/internal/ui/mounts/app/kv/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusOK, output),
	)
}

func newGetMntHandlerFailed() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("GET", "/v1/sys/internal/ui/mounts/app/kv/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusInternalServerError,
			&api.ErrorResponse{
				Errors: []string{"oh no!"},
			},
		),
	)
}

func newGetMntHandlerBadResponse() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("GET", "/v1/sys/internal/ui/mounts/app/kv/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusOK, "yahoo"),
	)
}

func newGetKVHandler() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("GET", "/v1/app/kv/data/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusOK, &api.Secret{
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"password": "swordfish",
				},
			},
		}),
	)
}

func newGetKVHandlerFailed() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("GET", "/v1/app/kv/data/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusInternalServerError,
			&api.ErrorResponse{
				Errors: []string{"oh no!"},
			},
		),
	)
}

func newGetKVHandlerBadResponse() http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest("GET", "/v1/app/kv/data/config"),
		VerifyHeaderKV("X-Vault-Token", "my-token"),
		RespondWithJSONEncoded(http.StatusOK, "yahoo"),
	)
}
