package ssm_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws/ssm"
	"github.com/phogolabs/cli/provider/aws/ssm/fake"
)

var _ = Describe("Provider", func() {
	var (
		flag *cli.StringFlag
		ctx  *cli.Context
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:  "listen-addr",
			Usage: "listen address of HTTP server",
			Metadata: cli.Map{
				"ssm_param": "/terraform/secret",
			},
		}

		ctx = &cli.Context{
			Command: &cli.Command{
				Name:  "app",
				Flags: []cli.Flag{flag},
			},
		}
	})

	Describe("SSM", func() {
		var (
			provider *ssm.Provider
			client   *fake.Client
		)

		BeforeEach(func() {
			client = &fake.Client{}
			client.GetReturns("swordfish", nil)

			provider = &ssm.Provider{
				Client: client,
			}
		})

		It("sets the value successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("swordfish"))

			Expect(client.GetCallCount()).To(Equal(1))
			Expect(client.GetArgsForCall(0)).To(Equal("/terraform/secret"))
		})

		Context("when the file path is not set", func() {
			BeforeEach(func() {
				flag.Metadata = cli.Map{}
			})

			It("does not set the value", func() {
				Expect(provider.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.GetReturns("", fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(provider.Provide(ctx)).To(MatchError("oh no!"))
			})
		})
	})
})
