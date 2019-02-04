package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("Flag", func() {
	Describe("StringFlag", func() {
		var flag *cli.StringFlag

		BeforeEach(func() {
			flag = &cli.StringFlag{
				Name:     "listen-addr",
				Value:    "9292",
				Usage:    "listen address of HTTP server",
				EnvVar:   "APP_LISTEN_ADDR",
				FilePath: "app.config",
			}
		})

		Describe("String", func() {
			It("returns the flag as string", func() {
				Expect(flag.String()).To(Equal("9292"))
			})
		})

		Describe("Set", func() {
			It("sets the value successfully", func() {
				Expect(flag.Set("8080")).To(Succeed())
				Expect(flag.Value).To(Equal("8080"))
			})
		})

		Describe("Get", func() {
			It("gets the value successfully", func() {
				Expect(flag.Get()).To(Equal("9292"))
			})
		})

		Describe("Validate", func() {
			It("validates the flag successfully", func() {
				Expect(flag.Validate()).To(Succeed())
			})

			Context("when the flag is required", func() {
				BeforeEach(func() {
					flag.Required = true
				})

				Context("when the flag's value is not set", func() {
					BeforeEach(func() {
						flag.Value = ""
					})

					It("returns an error", func() {
						Expect(flag.Validate()).To(MatchError("cli: flag -listen-addr is missing"))
					})
				})
			})
		})

		Describe("Definition", func() {
			It("returns the definition successfully", func() {
				definition := flag.Definition()
				Expect(definition).NotTo(BeNil())

				Expect(definition.Name).To(Equal(flag.Name))
				Expect(definition.Usage).To(Equal(flag.Usage))
				Expect(definition.FilePath).To(Equal(flag.FilePath))
				Expect(definition.EnvVar).To(Equal(flag.EnvVar))
				Expect(definition.Metadata).To(Equal(flag.Metadata))
			})
		})
	})
})
