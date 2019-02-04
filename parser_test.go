package cli_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/fake"
)

var _ = Describe("Parser", func() {
	var (
		flag       *fake.Flag
		definition *cli.FlagDefinition
		ctx        *cli.ParserContext
	)

	BeforeEach(func() {
		definition = &cli.FlagDefinition{
			Name:  "listen-addr",
			Usage: "listen address of HTTP server",
		}

		flag = &fake.Flag{}
		flag.DefinitionReturns(definition)

		ctx = &cli.ParserContext{
			Name:   "app",
			Flags:  []cli.Flag{flag},
			Output: GinkgoWriter,
		}
	})

	Describe("EnvParser", func() {
		var parser *cli.EnvParser

		BeforeEach(func() {
			definition.EnvVar = "APP_LISTEN_ADDR"

			parser = &cli.EnvParser{}
			Expect(os.Setenv(definition.EnvVar, "8080")).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.Unsetenv(definition.EnvVar)).To(Succeed())
		})

		It("sets the value from env variable", func() {
			Expect(parser.Parse(ctx)).To(Succeed())
			Expect(flag.SetCallCount()).To(Equal(1))
			Expect(flag.SetArgsForCall(0)).To(Equal("8080"))
		})

		Context("when the env var key is not set", func() {
			BeforeEach(func() {
				definition.EnvVar = ""
			})

			It("does not set the value", func() {
				Expect(parser.Parse(ctx)).To(Succeed())
				Expect(flag.SetCallCount()).To(BeZero())
			})
		})

		Context("when the env var is not set", func() {
			BeforeEach(func() {
				Expect(os.Setenv(definition.EnvVar, "")).To(Succeed())
			})

			It("does not set the value", func() {
				Expect(parser.Parse(ctx)).To(Succeed())
				Expect(flag.SetCallCount()).To(BeZero())
			})
		})

		Context("when setting the value fails", func() {
			BeforeEach(func() {
				flag.SetReturns(fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(parser.Parse(ctx)).To(MatchError("oh no!"))
			})
		})
	})

	Describe("FileParser", func() {
		var parser *cli.FileParser

		BeforeEach(func() {
			parser = &cli.FileParser{}

			tmpfile, err := ioutil.TempFile("", "example")
			Expect(err).To(BeNil())

			fmt.Fprint(tmpfile, "9292")

			definition.FilePath = tmpfile.Name()

			Expect(tmpfile.Close()).To(Succeed())
		})

		It("sets the value successfully", func() {
			Expect(parser.Parse(ctx)).To(Succeed())
			Expect(flag.SetCallCount()).To(Equal(1))
			Expect(flag.SetArgsForCall(0)).To(Equal("9292"))
		})

		Context("when the file does not exist", func() {
			BeforeEach(func() {
				definition.FilePath = "/tmp/file"
			})

			It("does not set the value", func() {
				Expect(parser.Parse(ctx)).To(Succeed())
				Expect(flag.SetCallCount()).To(BeZero())
			})
		})

		Context("when the file path is not set", func() {
			BeforeEach(func() {
				definition.FilePath = ""
			})

			It("does not set the value", func() {
				Expect(parser.Parse(ctx)).To(Succeed())
				Expect(flag.SetCallCount()).To(BeZero())
			})
		})

		Context("when setting the value fails", func() {
			BeforeEach(func() {
				flag.SetReturns(fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(parser.Parse(ctx)).To(MatchError("oh no!"))
			})
		})
	})

	Describe("FlagParser", func() {
		var parser *cli.FlagParser

		BeforeEach(func() {
			parser = &cli.FlagParser{}
			ctx.Args = []string{"-listen-addr=9292"}
		})

		It("sets the value successfully", func() {
			Expect(parser.Parse(ctx)).To(Succeed())
			Expect(flag.SetCallCount()).To(Equal(1))
			Expect(flag.SetArgsForCall(0)).To(Equal("9292"))
		})

		Context("when setting the value fails", func() {
			BeforeEach(func() {
				flag.SetReturns(fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(parser.Parse(ctx)).To(MatchError(`invalid value "9292" for flag -listen-addr: oh no!`))
			})
		})
	})
})
