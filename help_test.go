package cli_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/fake"
)

var _ = Describe("ShowHelp", func() {
	var (
		ctx    *cli.Context
		parent *cli.Context
		buffer *Buffer
	)

	BeforeEach(func() {
		buffer = NewBuffer()

		parent = &cli.Context{
			Manifest: &cli.Manifest{
				HideVersion: false,
				Version:     "1.0",
				Authors:     []*cli.Author{},
				Copyright:   "2020",
			},
			Writer: buffer,
			Command: &cli.Command{
				Name:      "root",
				HelpName:  "root_help",
				Usage:     "root_usage",
				UsageText: "root_usage_text",
				Commands: []*cli.Command{
					cli.NewHelpCommand(),
					&cli.Command{
						Name:      "action",
						HelpName:  "action_help",
						Usage:     "action_usage",
						UsageText: "action_usage_text",
					},
				},
			},
			Parent: nil,
		}
	})

	Context("when the command is main command", func() {
		It("shows the help for the command", func() {
			Expect(cli.NewHelpCommand().Action(parent)).To(Succeed())
			Expect(buffer).To(Say("COPYRIGHT"))
		})
	})

	Context("when the writer fails", func() {
		BeforeEach(func() {
			writer := &fake.Writer{}
			writer.WriteReturns(0, fmt.Errorf("oh no!"))
			parent.Writer = writer
		})

		It("returns an error", func() {
			Expect(cli.NewHelpCommand().Action(parent)).To(MatchError("oh no!"))
		})
	})

	Context("when help is executed for the command", func() {
		BeforeEach(func() {
			ctx = &cli.Context{
				Parent:  parent,
				Writer:  parent.Writer,
				Command: cli.NewHelpCommand(),
				Args:    []string{"action"},
			}
		})

		It("shows the help for the command", func() {
			Expect(cli.NewHelpCommand().Action(ctx)).To(Succeed())
			Expect(buffer).To(Say("action_help - action_usage"))
		})
	})

	Context("when help is executed for the command that does not exist", func() {
		BeforeEach(func() {
			ctx = &cli.Context{
				Parent:  parent,
				Writer:  parent.Writer,
				Command: cli.NewHelpCommand(),
				Args:    []string{"exec"},
			}
		})

		It("shows the help for the command", func() {
			Expect(cli.NewHelpCommand().Action(ctx)).To(Succeed())
			Expect(buffer).To(Say("No help topic for 'exec'"))
		})
	})

	Context("when the command is subcommand", func() {
		BeforeEach(func() {
			ctx = &cli.Context{
				Parent:  parent,
				Writer:  parent.Writer,
				Command: cli.NewHelpCommand(),
			}
		})

		It("shows the help for the command", func() {
			Expect(cli.NewHelpCommand().Action(ctx)).To(Succeed())
			Expect(buffer).To(Say("root_help - root_usage"))
		})
	})

})

var _ = Describe("ShowVersion", func() {
	It("shows the version successfully", func() {
		buffer := NewBuffer()

		ctx := &cli.Context{
			Parent: &cli.Context{
				Writer: buffer,
				Manifest: &cli.Manifest{
					Version: "BETA",
				},
				Command: &cli.Command{
					Name: "app",
				},
			},
		}

		Expect(cli.NewVersionCommand().Action(ctx)).To(Succeed())
		Expect(buffer).To(Say("app version BETA"))
	})
})
