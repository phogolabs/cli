package cli_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
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
		}

		ctx = &cli.Context{
			Command: &cli.Command{
				Name:  "app",
				Flags: []cli.Flag{flag},
			},
		}
	})

	Describe("EnvProvider", func() {
		var parser *cli.EnvProvider

		BeforeEach(func() {
			flag.EnvVar = "APP_LISTEN_ADDR"

			parser = &cli.EnvProvider{}
			Expect(os.Setenv(flag.EnvVar, "8080")).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.Unsetenv(flag.EnvVar)).To(Succeed())
		})

		It("sets the value from env variable", func() {
			Expect(parser.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("8080"))
		})

		Context("when the env var key is not set", func() {
			BeforeEach(func() {
				flag.EnvVar = ""
			})

			It("does not set the value", func() {
				Expect(parser.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeEmpty())
			})
		})

		Context("when the env var is not set", func() {
			BeforeEach(func() {
				Expect(os.Setenv(flag.EnvVar, "")).To(Succeed())
			})

			It("does not set the value", func() {
				Expect(parser.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeEmpty())
			})
		})

		Context("when setting the value fails", func() {
			var intFlag *cli.IntFlag

			BeforeEach(func() {
				intFlag = &cli.IntFlag{
					Name:   "num",
					Usage:  "number",
					EnvVar: "APP_NUM",
				}

				ctx = &cli.Context{
					Command: &cli.Command{
						Name:  "app",
						Flags: []cli.Flag{intFlag},
					},
				}

				Expect(os.Setenv(intFlag.EnvVar, "yep")).To(Succeed())
			})

			It("returns an error", func() {
				Expect(parser.Provide(ctx)).To(MatchError("strconv.ParseInt: parsing \"yep\": invalid syntax"))
			})
		})
	})

	Describe("FileProvider", func() {
		var parser *cli.FileProvider

		BeforeEach(func() {
			parser = &cli.FileProvider{}

			tmpfile, err := ioutil.TempFile("", "example")
			Expect(err).To(BeNil())

			fmt.Fprint(tmpfile, "9292")

			flag.FilePath = tmpfile.Name()

			Expect(tmpfile.Close()).To(Succeed())
		})

		It("sets the value successfully", func() {
			Expect(parser.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("9292"))
		})

		Context("when the file path is not valid", func() {
			It("returns an error", func() {
				flag.FilePath = "\\/"
				Expect(parser.Provide(ctx)).To(MatchError("syntax error in pattern"))
			})
		})

		Context("when the file does not exist", func() {
			BeforeEach(func() {
				flag.FilePath = "/tmp/file"
			})

			It("does not set the value", func() {
				Expect(parser.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})

		Context("when the file path is not set", func() {
			BeforeEach(func() {
				flag.FilePath = ""
			})

			It("does not set the value", func() {
				Expect(parser.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})

		Context("when setting the value fails", func() {
			var ip *cli.IPFlag

			BeforeEach(func() {
				ip = &cli.IPFlag{
					Name:     "listen-addr",
					FilePath: flag.FilePath,
				}

				ctx = &cli.Context{
					Command: &cli.Command{
						Name:  "app",
						Flags: []cli.Flag{ip},
					},
				}
			})

			It("returns an error", func() {
				Expect(parser.Provide(ctx)).To(MatchError("invalid IP Address: 9292"))
			})
		})
	})

	Describe("FlagProvider", func() {
		var parser *cli.FlagProvider

		BeforeEach(func() {
			parser = &cli.FlagProvider{}
			ctx.Args = []string{"-listen-addr=9292"}
		})

		It("sets the value successfully", func() {
			Expect(parser.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("9292"))
		})

		Context("when setting the value fails", func() {
			var ip *cli.IPFlag

			BeforeEach(func() {
				ip = &cli.IPFlag{
					Name: "listen-addr",
				}

				ctx = &cli.Context{
					Args: ctx.Args,
					Command: &cli.Command{
						Name:  "app",
						Flags: []cli.Flag{ip},
					},
				}
			})

			It("returns an error", func() {
				Expect(parser.Provide(ctx)).To(MatchError(`invalid value "9292" for flag -listen-addr: invalid IP Address: 9292`))
			})
		})
	})

	Describe("DefaultValueProvider", func() {
		var parser *cli.DefaultValueProvider

		BeforeEach(func() {
			flag.Value = "9292"
			parser = &cli.DefaultValueProvider{}
		})

		It("parses the flags successfully", func() {
			Expect(parser.Provide(ctx)).To(Succeed())
		})

		It("rollbacks the values successfully", func() {
			Expect(parser.Provide(ctx)).To(Succeed())

			flag.Value = "1010"

			Expect(parser.Rollback(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("9292"))
		})
	})
})
