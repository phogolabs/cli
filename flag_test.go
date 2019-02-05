package cli_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("Flag", func() {
	Describe("BoolFlag", func() {
		var flag *cli.BoolFlag

		BeforeEach(func() {
			flag = &cli.BoolFlag{
				Name:     "help",
				Usage:    "Show help info",
				EnvVar:   "APP_SHOW_HELP",
				FilePath: "help.txt",
			}
		})

		Describe("IsBoolFlag", func() {
			It("returns true", func() {
				Expect(flag.IsBoolFlag()).To(BeTrue())
			})
		})

		Describe("String", func() {
			It("returns the flag as string", func() {
				Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
			})
		})

		Describe("Set", func() {
			It("sets the value successfully", func() {
				Expect(flag.Set("true")).To(Succeed())
				Expect(flag.Value).To(BeTrue())
			})

			Context("when the value is not valid", func() {
				It("returns an error", func() {
					Expect(flag.Set("yahoo")).To(MatchError(`strconv.ParseBool: parsing "yahoo": invalid syntax`))
				})
			})
		})

		Describe("Get", func() {
			It("gets the value successfully", func() {
				Expect(flag.Get()).To(BeFalse())
			})
		})

		Describe("Validate", func() {
			It("validates the flag successfully", func() {
				Expect(flag.Validate()).To(Succeed())
			})

			Context("when the validation fails", func() {
				It("returns an error", func() {
					flag.ValidationFn = func(cli.Flag) error {
						return fmt.Errorf("oh no!")
					}

					Expect(flag.Validate()).To(MatchError("oh no!"))
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
				Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
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

			Context("when the validation fails", func() {
				It("returns an error", func() {
					flag.ValidationFn = func(cli.Flag) error {
						return fmt.Errorf("oh no!")
					}

					Expect(flag.Validate()).To(MatchError("oh no!"))
				})
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

	Describe("StringSliceFlag", func() {
		var flag *cli.StringSliceFlag

		BeforeEach(func() {
			flag = &cli.StringSliceFlag{
				Name:     "user",
				Value:    []string{"root"},
				Usage:    "List of all users",
				EnvVar:   "APP_LISTEN_ADDR",
				FilePath: "app.config",
			}
		})

		Describe("String", func() {
			It("returns the flag as string", func() {
				Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
			})
		})

		Describe("Set", func() {
			It("sets the value successfully", func() {
				Expect(flag.Set("guest, admin")).To(Succeed())
				Expect(flag.Value).To(HaveLen(2))
				Expect(flag.Value).To(ContainElement("guest"))
				Expect(flag.Value).To(ContainElement("admin"))
			})
		})

		Describe("Get", func() {
			It("gets the value successfully", func() {
				Expect(flag.Get()).To(Equal(flag.Value))
			})
		})

		Describe("Validate", func() {
			It("validates the flag successfully", func() {
				Expect(flag.Validate()).To(Succeed())
			})

			Context("when the validation fails", func() {
				It("returns an error", func() {
					flag.ValidationFn = func(cli.Flag) error {
						return fmt.Errorf("oh no!")
					}

					Expect(flag.Validate()).To(MatchError("oh no!"))
				})
			})

			Context("when the flag is required", func() {
				BeforeEach(func() {
					flag.Required = true
				})

				Context("when the flag's value is not set", func() {
					BeforeEach(func() {
						flag.Value = nil
					})

					It("returns an error", func() {
						Expect(flag.Validate()).To(MatchError("cli: flag -user is missing"))
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
