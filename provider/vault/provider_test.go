package vault_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/vault"
	"github.com/phogolabs/vault/fake"
)

var _ = Describe("Provider", func() {
	var (
		provider *vault.Provider
		server   *Server
		ctx      *cli.Context
		handlers []http.HandlerFunc
	)

	BeforeEach(func() {
		handlers = []http.HandlerFunc{
			newAuthHandler(),
			newGetMntHandler(),
			newGetKVHandler(),
		}

		provider = &vault.Provider{}
	})

	JustBeforeEach(func() {
		server = NewServer()
		server.AppendHandlers(handlers...)

		ctx = &cli.Context{
			Command: &cli.Command{
				Name: "app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "vault-addr",
						Value: server.URL(),
					},
					&cli.StringFlag{
						Name:  "vault-auth-mount-path",
						Value: "kubernetes",
					},
					&cli.StringFlag{
						Name:  "vault-auth-role",
						Value: "admin",
					},
					&cli.StringFlag{
						Name:  "vault-auth-kube-jwt",
						Value: "kubo",
					},
					&cli.StringFlag{
						Name: "password",
						Converter: &cli.JSONPathConverter{
							Path: "$.password",
						},
						Metadata: map[string]string{
							"vault_path": "/app/kv/config",
						},
					},
				},
			},
		}
	})

	AfterEach(func() {
		Expect(provider.Rollback(ctx)).To(Succeed())

		server.Close()
	})

	It("parses the flags successfully", func() {
		Expect(provider.Provide(ctx)).To(Succeed())
		Expect(ctx.String("password")).To(Equal("swordfish"))
	})

	Context("when the token is provided", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newGetMntHandler(),
				newGetKVHandler(),
			}
		})

		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "vault-addr",
							Value: server.URL(),
						},
						&cli.StringFlag{
							Name:  "vault-token",
							Value: "my-token",
						},
						&cli.StringFlag{
							Name: "password",
							Converter: &cli.JSONPathConverter{
								Path: "$.password",
							},
							Metadata: map[string]string{
								"vault_path": "/app/kv/config",
							},
						},
					},
				},
			}
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(ctx.String("password")).To(Equal("swordfish"))
		})
	})

	Context("when the authentication fails", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandlerFailed(),
				newAuthHandlerFailed(),
				newAuthHandlerFailed(),
			}
		})

		JustBeforeEach(func() {
			server.SetAllowUnhandledRequests(true)
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx).Error()).To(ContainSubstring("Code: 500"))
		})
	})

	Context("when the fetcher fails", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandler(),
				newGetMntHandlerFailed(),
				newGetMntHandlerFailed(),
				newGetMntHandlerFailed(),
			}
		})

		JustBeforeEach(func() {
			server.SetAllowUnhandledRequests(true)
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx).Error()).To(ContainSubstring("Code: 500"))
		})
	})

	Context("when setting the flag fails", func() {
		JustBeforeEach(func() {
			flags := ctx.Command.Flags
			flags[len(flags)-1] = &cli.IntFlag{
				Name: "password",
				Converter: &cli.JSONPathConverter{
					Path: "$.password",
				},
				Metadata: map[string]string{
					"vault_path": "/app/kv/config",
				},
			}
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx)).To(MatchError("strconv.ParseInt: parsing \"swordfish\": invalid syntax"))
		})
	})

	Context("when the json path is not valid", func() {
		JustBeforeEach(func() {
			flags := ctx.Command.Flags
			flags[len(flags)-1] = &cli.StringFlag{
				Name: "password",
				Converter: &cli.JSONPathConverter{
					Path: "$.$",
				},
				Metadata: map[string]string{
					"vault_path": "/app/kv/config",
				},
			}
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx)).To(MatchError("expression don't support in filter"))
		})
	})

	Context("when the repository is already initialized", func() {
		var repository *fake.Repository

		BeforeEach(func() {
			repository = &fake.Repository{}
			repository.SecretReturns(map[string]interface{}{
				"password": "swordfish",
			}, nil)

			provider.Repository = repository
		})

		AfterEach(func() {
			provider.Repository = nil
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(ctx.String("password")).To(Equal("swordfish"))
		})
	})

	Context("when the fetcher cannot be initialized", func() {
		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "password",
							Converter: &cli.JSONPathConverter{
								Path: "$.password",
							},
							Metadata: map[string]string{
								"vault_path": "/app/kv/config",
							},
						},
					},
				},
			}
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(provider.Repository).To(BeNil())
			Expect(ctx.String("password")).To(BeEmpty())
		})
	})

	Context("when the valut-addr is not valid", func() {
		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "vault-addr",
							Value: "://address",
						},
						&cli.StringFlag{
							Name: "password",
							Metadata: map[string]string{
								"vault_path": "/app/kv/config",
							},
						},
					},
				},
			}
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx)).To(MatchError("parse ://address: missing protocol scheme"))
		})
	})
})
