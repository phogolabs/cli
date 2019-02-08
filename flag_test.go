package cli_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("BoolFlag", func() {
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

var _ = Describe("StringFlag", func() {
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

var _ = Describe("StringSliceFlag", func() {
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

var _ = Describe("URLFlag", func() {
	var flag *cli.URLFlag

	BeforeEach(func() {
		value, err := url.Parse("http://example.com")
		Expect(err).To(BeNil())

		flag = &cli.URLFlag{
			Name:     "listen-addr",
			Value:    value,
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
			Expect(flag.Set("http://google.com")).To(Succeed())
			Expect(flag.Value.String()).To(Equal("http://google.com"))
		})
	})

	Context("when the value cannot be parsed", func() {
		It("returns an error", func() {
			Expect(flag.Set("://wrong")).To(MatchError("parse ://wrong: missing protocol scheme"))
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

var _ = Describe("JSONFlag", func() {
	var flag *cli.JSONFlag

	BeforeEach(func() {
		flag = &cli.JSONFlag{
			Name: "map",
			Value: map[string]interface{}{
				"id":   0,
				"user": "root",
			},
			EnvVar: "APP_MAP",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(ContainSubstring("id => 0"))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			m := map[string]string{
				"key": "value",
			}

			data, err := json.Marshal(&m)
			Expect(err).To(BeNil())

			Expect(flag.Set(string(data))).To(Succeed())
			Expect(flag.Value).To(HaveKeyWithValue("key", "value"))
		})
	})

	Context("when the value cannot be parsed", func() {
		It("returns an error", func() {
			Expect(flag.Set("wrong")).To(MatchError("invalid character 'w' looking for beginning of value"))
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
					Expect(flag.Validate()).To(MatchError("cli: flag -map is missing"))
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

var _ = Describe("YAMLFlag", func() {
	var flag *cli.YAMLFlag

	BeforeEach(func() {
		flag = &cli.YAMLFlag{
			Name: "map",
			Value: map[string]interface{}{
				"id":   0,
				"user": "root",
			},
			EnvVar: "APP_MAP",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(ContainSubstring("id => 0"))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			m := map[string]string{
				"key": "value",
			}

			data, err := json.Marshal(&m)
			Expect(err).To(BeNil())

			Expect(flag.Set(string(data))).To(Succeed())
			Expect(flag.Value).To(HaveKeyWithValue("key", "value"))
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
					Expect(flag.Validate()).To(MatchError("cli: flag -map is missing"))
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

var _ = Describe("TimeFlag", func() {
	var flag *cli.TimeFlag

	BeforeEach(func() {
		flag = &cli.TimeFlag{
			Name:   "time",
			Value:  time.Now(),
			EnvVar: "APP_TIME",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			t := time.Now()
			value := t.Format(time.UnixDate)
			Expect(flag.Set(value)).To(Succeed())
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
					flag.Value = time.Time{}
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -time is missing"))
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

var _ = Describe("DurationFlag", func() {
	var flag *cli.DurationFlag

	BeforeEach(func() {
		flag = &cli.DurationFlag{
			Name:   "time",
			Value:  time.Second,
			EnvVar: "APP_TIME",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("10s")).To(Succeed())
			Expect(flag.Value).To(Equal(10 * time.Second))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -time is missing"))
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

var _ = Describe("IntFlag", func() {
	var flag *cli.IntFlag

	BeforeEach(func() {
		flag = &cli.IntFlag{
			Name:  "num",
			Value: 66,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("99")).To(Succeed())
			Expect(flag.Value).To(Equal(99))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(66))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -num is missing"))
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

var _ = Describe("Int64Flag", func() {
	var flag *cli.Int64Flag

	BeforeEach(func() {
		flag = &cli.Int64Flag{
			Name:  "num",
			Value: 66,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("99")).To(Succeed())
			Expect(flag.Value).To(Equal(int64(99)))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(int64(66)))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -num is missing"))
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

var _ = Describe("UIntFlag", func() {
	var flag *cli.UIntFlag

	BeforeEach(func() {
		flag = &cli.UIntFlag{
			Name:  "num",
			Value: 66,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("99")).To(Succeed())
			Expect(flag.Value).To(Equal(uint(99)))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(uint(66)))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -num is missing"))
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

var _ = Describe("UInt64Flag", func() {
	var flag *cli.UInt64Flag

	BeforeEach(func() {
		flag = &cli.UInt64Flag{
			Name:  "num",
			Value: 66,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("99")).To(Succeed())
			Expect(flag.Value).To(Equal(uint64(99)))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(uint64(66)))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -num is missing"))
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

var _ = Describe("Float64Flag", func() {
	var flag *cli.Float64Flag

	BeforeEach(func() {
		flag = &cli.Float64Flag{
			Name:  "num",
			Value: 66,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("99")).To(Succeed())
			Expect(flag.Value).To(Equal(float64(99)))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(float64(66)))
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
					flag.Value = 0
				})

				It("returns an error", func() {
					Expect(flag.Validate()).To(MatchError("cli: flag -num is missing"))
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
