package cli_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/url"
	"sort"
	"time"

	"github.com/phogolabs/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BoolFlag", func() {
	var flag *cli.BoolFlag

	BeforeEach(func() {
		flag = &cli.BoolFlag{
			Name:   "help",
			Path:   "help.txt",
			Usage:  "Show help info",
			EnvVar: "APP_SHOW_HELP",
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
})

var _ = Describe("StringFlag", func() {
	var flag *cli.StringFlag

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:   "listen-addr",
			Usage:  "listen address of HTTP server",
			EnvVar: "APP_LISTEN_ADDR",
			Path:   "app.config",
			Value:  ":9292",
		}
	})

	Describe("String", func() {
		It("returns the flag as string", func() {
			Expect(flag.String()).To(Equal(cli.FlagFormat(flag)))
		})
	})

	Describe("Set", func() {
		It("sets the value successfully", func() {
			Expect(flag.Set(":8080")).To(Succeed())
			Expect(flag.Value).To(Equal(":8080"))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(":9292"))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'listen-addr' not found"))
				})
			})
		})
	})
})

var _ = Describe("StringSliceFlag", func() {
	var flag *cli.StringSliceFlag

	BeforeEach(func() {
		flag = &cli.StringSliceFlag{
			Name:   "user",
			Path:   "app.config",
			Usage:  "List of all users",
			EnvVar: "APP_LISTEN_ADDR",
			Value:  []string{"root"},
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(ctx *cli.Context, v interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'user' not found"))
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
			Name:   "listen-addr",
			Value:  value,
			Usage:  "listen address of HTTP server",
			EnvVar: "APP_LISTEN_ADDR",
			Path:   "app.config",
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
			Expect(flag.Set("://wrong")).To(HaveOccurred())
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(flag.Value))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'listen-addr' not found"))
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

				_, err = flag.ReadFrom(bytes.NewBuffer(data))
				Expect(err).NotTo(HaveOccurred())

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

				_, err = flag.ReadFrom(bytes.NewBuffer(data))
				Expect(err).NotTo(HaveOccurred())

				user, ok := flag.Value.(*User)
				Expect(ok).To(BeTrue())
				Expect(user.Name).To(Equal("John"))
			})
		})
	})

	Context("when the value cannot be parsed", func() {
		It("returns an error", func() {
			_, err := flag.ReadFrom(bytes.NewBufferString("wrong"))
			Expect(err).To(MatchError("invalid character 'w' looking for beginning of value"))
		})
	})

	Describe("Get", func() {
		It("gets the value successfully", func() {
			Expect(flag.Get()).To(Equal(flag.Path))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(ctx *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'map' not found"))
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

				_, err = flag.ReadFrom(bytes.NewBuffer(data))
				Expect(err).NotTo(HaveOccurred())

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
			Expect(flag.Get()).To(Equal(flag.Path))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(ctx *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'map' not found"))
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
			Expect(flag.Get()).To(Equal(flag.Path))
		})
	})

	Describe("Validate", func() {
		It("validates the flag successfully", func() {
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'map' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'time' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'time' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'num' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'ip' not found"))
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
			Expect(flag.Validate(&cli.Context{})).To(Succeed())
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				flag.Validator = cli.ValidatorFunc(func(_ *cli.Context, _ interface{}) error {
					return fmt.Errorf("oh no!")
				})

				Expect(flag.Validate(&cli.Context{})).To(MatchError("oh no!"))
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
					Expect(flag.Validate(&cli.Context{})).To(MatchError("flag 'mac' not found"))
				})
			})
		})
	})
})

var _ = Describe("FlagsByName", func() {
	It("sorts the flags correctly", func() {
		var (
			alpha = &cli.StringFlag{Name: "alpha"}
			beta  = &cli.StringFlag{Name: "beta"}
			flags = cli.FlagsByName{beta, alpha}
		)

		sort.Sort(flags)

		Expect(flags[0]).To(Equal(alpha))
		Expect(flags[1]).To(Equal(beta))
	})
})

var _ = Describe("FlagAccessor", func() {
	var (
		accessor *cli.FlagAccessor
		flag     *cli.StringFlag
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:   "listen-addr",
			Path:   "app.config",
			Usage:  "listen address of HTTP server",
			EnvVar: "APP_LISTEN_ADDR",
			Value:  "9292",
		}

		accessor = &cli.FlagAccessor{
			Flag: flag,
		}
	})

	It("returns the definition successfully", func() {
		Expect(accessor.Name()).To(Equal(flag.Name))
		Expect(accessor.Path()).To(Equal(flag.Path))
		Expect(accessor.Usage()).To(Equal(flag.Usage))
		Expect(accessor.EnvVar()).To(Equal(flag.EnvVar))
		Expect(accessor.Value()).To(Equal(flag.Value))
	})
})
