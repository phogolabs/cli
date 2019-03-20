package cli_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
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
				flag.Validator = cli.ValidatorFunc(func(v interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate()).To(MatchError("oh no!"))
			})
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
				flag.Validator = cli.ValidatorFunc(func(v interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
			Expect(flag.Set("guest")).To(Succeed())
			Expect(flag.Set("admin")).To(Succeed())
			Expect(flag.Set("usr1,usr2")).To(Succeed())

			Expect(flag.Value).To(HaveLen(4))
			Expect(flag.Value).To(ContainElement("guest"))
			Expect(flag.Value).To(ContainElement("admin"))
			Expect(flag.Value).To(ContainElement("usr1,usr2"))
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
				flag.Validator = cli.ValidatorFunc(func(v interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
})

var _ = Describe("JSONFlag", func() {
	var flag *cli.JSONFlag

	BeforeEach(func() {
		value := &map[string]interface{}{
			"id":   0,
			"user": "root",
		}

		flag = &cli.JSONFlag{
			Name:   "map",
			Value:  value,
			EnvVar: "APP_MAP",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(ContainSubstring("id => 0"))
		})
	})

	Describe("Set", func() {
		ItSetsTheValue := func() {
			It("sets the value successfully", func() {
				m := map[string]string{
					"key": "value",
				}

				data, err := json.Marshal(&m)
				Expect(err).To(BeNil())
				Expect(flag.Set(fmt.Sprintf("'%v'", string(data)))).To(Succeed())

				value, ok := flag.Value.(*map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(*value).To(HaveKeyWithValue("key", "value"))
			})
		}

		ItSetsTheValue()

		Context("when the value is not set", func() {
			BeforeEach(func() {
				flag.Value = nil
			})
			ItSetsTheValue()
		})

		Context("when the value is map of map", func() {
			It("sets the value successfully", func() {
				type User struct {
					Name string `json:"name"`
				}

				flag.Value = &User{}

				m := map[string]string{
					"name": "John",
				}

				data, err := json.Marshal(&m)
				Expect(err).To(BeNil())

				Expect(flag.Set(string(data))).To(Succeed())

				user, ok := flag.Value.(*User)
				Expect(ok).To(BeTrue())
				Expect(user.Name).To(Equal("John"))
			})
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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
})

var _ = Describe("YAMLFlag", func() {
	var flag *cli.YAMLFlag

	BeforeEach(func() {
		flag = &cli.YAMLFlag{
			Name: "map",
			Value: &map[string]interface{}{
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
		ItSetsTheValue := func() {
			It("sets the value successfully", func() {
				m := map[string]string{
					"key": "value",
				}

				data, err := json.Marshal(&m)
				Expect(err).To(BeNil())

				Expect(flag.Set(string(data))).To(Succeed())

				value, ok := flag.Value.(*map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(*value).To(HaveKeyWithValue("key", "value"))
			})
		}

		ItSetsTheValue()

		Context("when the value is not set", func() {
			BeforeEach(func() {
				flag.Value = nil
			})
			ItSetsTheValue()

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
})

var _ = Describe("XMLFlag", func() {
	var flag *cli.XMLFlag

	type T struct {
		ID   int
		Name string
	}

	BeforeEach(func() {
		flag = &cli.XMLFlag{
			Name: "map",
			Value: &T{
				ID:   0,
				Name: "root",
			},
			EnvVar: "APP_MAP",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(ContainSubstring("--map value"))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			m := T{
				ID:   12345,
				Name: "guest",
			}

			data, err := xml.Marshal(&m)
			Expect(err).To(BeNil())

			Expect(flag.Set(string(data))).To(Succeed())
		})

		Context("when the value is not set", func() {
			It("sets the value successfully", func() {
				m := T{
					ID:   12345,
					Name: "guest",
				}

				data, err := xml.Marshal(&m)
				Expect(err).To(BeNil())

				flag.Value = nil
				Expect(flag.Set(string(data))).To(Succeed())
				Expect(flag.Value).To(BeNil())
			})
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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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

		Context("when the value is invalid", func() {
			It("returns an error", func() {
				Expect(flag.Set("yahoo")).To(MatchError(`strconv.ParseInt: parsing "yahoo": invalid syntax`))
			})
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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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

		Context("when the value is invalid", func() {
			It("returns an error", func() {
				Expect(flag.Set("yahoo")).To(MatchError(`strconv.ParseUint: parsing "yahoo": invalid syntax`))
			})
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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
})

var _ = Describe("Float32Flag", func() {
	var flag *cli.Float32Flag

	BeforeEach(func() {
		flag = &cli.Float32Flag{
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
			Expect(flag.Value).To(Equal(float32(99)))
		})

		Context("when the value is invalid", func() {
			It("sets the value successfully", func() {
				Expect(flag.Set("yahoo").Error()).To(Equal(`strconv.ParseFloat: parsing "yahoo": invalid syntax`))
			})
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(float32(66)))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate()).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
})

var _ = Describe("IPFlag", func() {
	var flag *cli.IPFlag

	BeforeEach(func() {
		flag = &cli.IPFlag{
			Name:  "ip",
			Value: net.ParseIP("127.0.0.1"),
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("127.0.1.1")).To(Succeed())
			Expect(flag.Value).To(Equal(net.ParseIP("127.0.1.1")))
		})

		Context("when the value is invalid", func() {
			It("returns an error", func() {
				Expect(flag.Set("yahoo")).To(MatchError("invalid IP Address: yahoo"))
			})
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(net.ParseIP("127.0.0.1")))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate()).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
					Expect(flag.Validate()).To(MatchError("cli: flag -ip is missing"))
				})
			})
		})
	})
})

var _ = Describe("HardwareAddrFlag", func() {
	var flag *cli.HardwareAddrFlag

	BeforeEach(func() {
		mac, err := net.ParseMAC("01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00")
		Expect(err).To(BeNil())

		flag = &cli.HardwareAddrFlag{
			Name:  "mac",
			Value: mac,
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set("01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00")).To(Succeed())
		})

		Context("when the value is not valid", func() {
			It("returns an error", func() {
				Expect(flag.Set("yahoo")).To(MatchError("address yahoo: invalid MAC address"))
			})
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
				flag.Validator = cli.ValidatorFunc(func(_ interface{}) error {
					return fmt.Errorf("oh no!")
				})

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
					Expect(flag.Validate()).To(MatchError("cli: flag -mac is missing"))
				})
			})
		})
	})
})

var _ = Describe("FlagAccessor", func() {
	var (
		accessor *cli.FlagAccessor
		flag     *cli.StringFlag
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:     "listen-addr",
			Value:    "9292",
			Usage:    "listen address of HTTP server",
			EnvVar:   "APP_LISTEN_ADDR",
			FilePath: "app.config",
			Metadata: map[string]string{
				"key": "meta",
			},
		}

		accessor = &cli.FlagAccessor{
			Flag: flag,
		}
	})

	It("returns the definition successfully", func() {
		Expect(accessor.Name()).To(Equal(flag.Name))
		Expect(accessor.Usage()).To(Equal(flag.Usage))
		Expect(accessor.FilePath()).To(Equal(flag.FilePath))
		Expect(accessor.EnvVar()).To(Equal(flag.EnvVar))
		Expect(accessor.Metadata()).To(Equal(flag.Metadata))
		Expect(accessor.Value()).To(Equal(flag.Value))
		Expect(accessor.MetaKey("key")).To(Equal("meta"))
	})

	Describe("SetValue", func() {
		It("sets the value", func() {
			Expect(accessor.SetValue("1212")).To(Succeed())
			Expect(flag.Value).To(Equal("1212"))
		})

		Context("when the value is not compatible", func() {
			It("returns an error", func() {
				Expect(accessor.SetValue(1)).To(MatchError("reflect.Set: value of type int is not assignable to type string"))
			})
		})

		Context("when the converter returns an error", func() {
			BeforeEach(func() {
				flag.Converter = cli.ConverterFunc(func(_ interface{}) (interface{}, error) {
					return nil, fmt.Errorf("oh no!")
				})
			})

			It("returns an error", func() {
				Expect(accessor.SetValue(1)).To(MatchError("oh no!"))
			})
		})
	})
})

var _ = Describe("JSONPath", func() {
	var (
		converter cli.JSONPath
		value     map[string]string
	)

	BeforeEach(func() {
		value = map[string]string{
			"password": "swordfish",
		}
		converter = cli.JSONPath("$.password")
	})

	It("converts the value successfully", func() {
		v, err := converter.Convert(value)
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("swordfish"))
	})

	Context("when the expression is wrong", func() {
		BeforeEach(func() {
			converter = cli.JSONPath("$.$")
		})

		It("returns an error", func() {
			v, err := converter.Convert(value)
			Expect(err).To(MatchError("expression don't support in filter"))
			Expect(v).To(BeNil())
		})
	})
})
