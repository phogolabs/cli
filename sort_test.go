package cli_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("Sort", func() {
	Describe("FlagsByName", func() {
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

	Describe("CommandsByName", func() {
		It("sorts the commands correctly", func() {
			var (
				alpha    = &cli.Command{Name: "alpha"}
				beta     = &cli.Command{Name: "beta"}
				commands = cli.CommandsByName{beta, alpha}
			)

			sort.Sort(commands)

			Expect(commands[0]).To(Equal(alpha))
			Expect(commands[1]).To(Equal(beta))
		})
	})
})
