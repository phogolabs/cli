package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/fake"
)

var _ = Describe("Format", func() {
	var (
		flag       *fake.Flag
		definition *cli.FlagDefinition
	)

	BeforeEach(func() {
		definition = &cli.FlagDefinition{
			Name:     "log-level, l",
			Usage:    "Application log level",
			EnvVar:   "LOG_LEVEL,LOG_LVL",
			FilePath: "my-config.cfg",
		}

		flag = &fake.Flag{}
		flag.GetReturns("info")
		flag.DefinitionReturns(definition)
	})

	It("formats a flag successfully", func() {
		help := cli.FlagFormat(flag)
		Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
	})

	Context("when the FilePath is not set", func() {
		BeforeEach(func() {
			definition.FilePath = ""
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [$LOG_LEVEL, $LOG_LVL]"))
		})
	})

	Context("when the EnvVar is not set", func() {
		BeforeEach(func() {
			definition.EnvVar = ""
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [my-config.cfg]"))
		})
	})

	Context("when the value is boolean", func() {
		BeforeEach(func() {
			flag.GetReturns(false)
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level, -l\tApplication log level [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
		})
	})

	Context("when the value is nil", func() {
		BeforeEach(func() {
			flag.GetReturns(nil)
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
		})
	})
})
