package cli_test

import (
	"github.com/phogolabs/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Format", func() {
	var (
		flag *cli.StringFlag
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:   "log-level, l",
			Path:   "logger.conf",
			Usage:  "Application log level",
			EnvVar: "LOG_LEVEL, LOG_LVL",
			Value:  "info",
		}
	})

	It("formats a flag successfully", func() {
		help := cli.FlagFormat(flag)
		Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [$LOG_LEVEL, $LOG_LVL] [logger.conf]"))
	})

	Context("when the Path is not set", func() {
		BeforeEach(func() {
			flag.Path = ""
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
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level (default: info) [logger.conf]"))
		})
	})

	Context("when the value is boolean", func() {
		var flag *cli.BoolFlag

		BeforeEach(func() {
			flag = &cli.BoolFlag{
				Name:   "log-level, l",
				Path:   "logger.conf",
				Usage:  "Application log level",
				EnvVar: "LOG_LEVEL, LOG_LVL",
				Value:  true,
			}
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level, -l\tApplication log level [$LOG_LEVEL, $LOG_LVL] [logger.conf]"))
		})
	})

	Context("when the value is nil", func() {
		var flag *cli.YAMLFlag

		BeforeEach(func() {
			flag = &cli.YAMLFlag{
				Name:   "log-level, l",
				Path:   "logger.conf",
				Usage:  "Application log level",
				EnvVar: "LOG_LEVEL, LOG_LVL",
			}
		})

		It("formats a flag successfully", func() {
			help := cli.FlagFormat(flag)
			Expect(help).To(Equal("--log-level value, -l value\tApplication log level [$LOG_LEVEL, $LOG_LVL] [logger.conf]"))
		})
	})
})
