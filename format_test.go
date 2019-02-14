package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("Format", func() {
	var (
		flag *cli.StringFlag
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:     "log-level, l",
			Usage:    "Application log level",
			EnvVar:   "LOG_LEVEL,LOG_LVL",
			FilePath: "my-config.cfg",
			Value:    "info",
		}
	})

	It("formats a flag successfully", func() {
		help := cli.FlagFormat(flag)
		Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
	})

	Context("when the FilePath is not set", func() {
		BeforeEach(func() {
			flag.FilePath = ""
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [$LOG_LEVEL, $LOG_LVL]"))
		})
	})

	Context("when the EnvVar is not set", func() {
		BeforeEach(func() {
			flag.EnvVar = ""
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [my-config.cfg]"))
		})
	})

	Context("when the value is boolean", func() {
		var flag *cli.BoolFlag

		BeforeEach(func() {
			flag = &cli.BoolFlag{
				Name:     "log-level, l",
				Usage:    "Application log level",
				EnvVar:   "LOG_LEVEL,LOG_LVL",
				FilePath: "my-config.cfg",
				Value:    true,
			}
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level, -l\tApplication log level [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
		})
	})

	Context("when the value is nil", func() {
		var flag *cli.YAMLFlag

		BeforeEach(func() {
			flag = &cli.YAMLFlag{
				Name:     "log-level, l",
				Usage:    "Application log level",
				EnvVar:   "LOG_LEVEL,LOG_LVL",
				FilePath: "my-config.cfg",
			}
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level [$LOG_LEVEL, $LOG_LVL] [my-config.cfg]"))
		})
	})
})
