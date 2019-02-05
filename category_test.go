package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("Category", func() {
	var category *cli.CommandCategory

	BeforeEach(func() {
		category = &cli.CommandCategory{
			Name: "main",
			Commands: []*cli.Command{
				&cli.Command{
					Name: "cmd1",
				},
				&cli.Command{
					Name:   "cmd2",
					Hidden: true,
				},
				&cli.Command{
					Name: "cmd3",
				},
			},
		}
	})

	Describe("VisibleCommands", func() {
		It("returns visible commands", func() {
			cmds := category.VisibleCommands()
			Expect(cmds).To(HaveLen(2))
			Expect(cmds[0].Name).To(Equal("cmd1"))
			Expect(cmds[1].Name).To(Equal("cmd3"))
		})
	})

})
