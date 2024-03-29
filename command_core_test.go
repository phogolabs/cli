package cli_test

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/phogolabs/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command", func() {
	var cmd *cli.Command

	BeforeEach(func() {
		cmd = &cli.Command{
			Name:      "run",
			Aliases:   []string{"start", "exec"},
			Usage:     "run_usage",
			UsageText: "run_usage_text",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "dir, d",
				},
				&cli.StringFlag{
					Name:   "path, p",
					Hidden: true,
				},
			},
			Commands: []*cli.Command{
				&cli.Command{
					Name:     "child1",
					Category: "cat1",
				},
				&cli.Command{
					Name:      "child2",
					Hidden:    true,
					Category:  "cat1",
					Usage:     "child2_usage",
					UsageText: "child2_usage_text",
				},
				&cli.Command{
					Name:     "child3",
					Category: "cat2",
				},
			},
		}
	})

	Describe("RunWithContext", func() {
		var (
			ctx    *cli.Context
			buffer *bytes.Buffer
		)

		BeforeEach(func() {
			buffer = &bytes.Buffer{}

			cmd.Action = func(ctx *cli.Context) error {
				Expect(ctx.Command).To(Equal(cmd))
				return nil
			}

			ctx = &cli.Context{
				Writer:  buffer,
				Command: cmd,
			}
		})

		It("runs the commandf successfully", func() {
			Expect(cmd.RunWithContext(ctx)).To(Succeed())
		})

		Context("when the parser fails", func() {
			BeforeEach(func() {
				ctx.Args = []string{"-unknonw-flag"}

				cmd.OnUsageError = func(ctx *cli.Context, err error) error {
					Expect(ctx.Command).To(Equal(cmd))
					Expect(err).To(MatchError("flag provided but not defined: -unknonw-flag"))
					return err
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("flag provided but not defined: -unknonw-flag"))
			})

			Context("when we hide the error", func() {
				BeforeEach(func() {
					cmd.OnUsageError = func(ctx *cli.Context, err error) error {
						Expect(ctx.Command).To(Equal(cmd))
						Expect(err).To(MatchError("flag provided but not defined: -unknonw-flag"))
						return nil
					}
				})

				It("does not return an error", func() {
					Expect(cmd.RunWithContext(ctx)).To(Succeed())
				})
			})
		})

		Context("when the flag validation fails", func() {
			BeforeEach(func() {
				flag := &cli.StringFlag{
					Name:     "name, n",
					Required: true,
				}

				cmd.Flags = append(cmd.Flags, flag)
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("flag 'name, n' not found"))
			})
		})

		Context("when a subcommand is executed", func() {
			BeforeEach(func() {
				ctx.Command.Commands[0].Action = func(child *cli.Context) error {
					Expect(child.Command).To(Equal(ctx.Command.Commands[0]))
					return nil
				}

				ctx.Args = []string{"child1"}
			})

			It("runs the command successfully", func() {
				Expect(cmd.RunWithContext(ctx)).To(Succeed())
			})

			Context("when the subcommand is not found", func() {
				BeforeEach(func() {
					ctx.Args = []string{"child69"}
				})

				It("shows the help", func() {
					Expect(cmd.RunWithContext(ctx)).To(Succeed())
					Expect(buffer.String()).To(Equal("No help topic for 'child69'\n"))
				})
			})

			Context("when the -h flag is passed", func() {
				BeforeEach(func() {
					ctx.Args = []string{"-h"}
				})

				It("shows the help", func() {
					Expect(cmd.RunWithContext(ctx)).To(Succeed())
					Expect(buffer.String()).To(ContainSubstring("run - run_usage"))
				})
			})

			Context("when the handler is not set", func() {
				BeforeEach(func() {
					ctx.Args = []string{"child2"}
				})

				It("shows the help", func() {
					Expect(cmd.RunWithContext(ctx)).To(Succeed())
					Expect(buffer.String()).To(ContainSubstring("run child2 - child2_usage"))
				})
			})
		})

		Context("when the before handler returns an error", func() {
			BeforeEach(func() {
				cmd.Before = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("oh no!"))
			})
		})

		Context("when the before init handler returns an error", func() {
			BeforeEach(func() {
				cmd.BeforeInit = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("oh no!"))
			})
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				cmd.Action = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("oh no!"))
			})
		})

		Context("when the after handler returns an error", func() {
			BeforeEach(func() {
				cmd.After = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("oh no!"))
			})
		})

		Context("when the after init handler returns an error", func() {
			BeforeEach(func() {
				cmd.AfterInit = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}
			})

			It("returns an error", func() {
				Expect(cmd.RunWithContext(ctx)).To(MatchError("oh no!"))
			})
		})
	})

	Describe("Names", func() {
		It("returns all names", func() {
			names := cmd.Names()
			Expect(names).To(HaveLen(3))
			Expect(names).To(ContainElement("run"))
			Expect(names).To(ContainElement("start"))
			Expect(names).To(ContainElement("exec"))
		})
	})

	Describe("VisibleFlags", func() {
		It("returns the visible flags", func() {
			ctx := &cli.Context{
				Writer:  GinkgoWriter,
				Command: cmd,
			}

			Expect(cmd.RunWithContext(ctx)).To(Succeed())

			flags := cmd.VisibleFlags()
			Expect(flags).To(HaveLen(2))
		})
	})

	Describe("VisibleCommands", func() {
		It("returns the visible commands", func() {
			cmds := cmd.VisibleCommands()
			Expect(cmds).To(HaveLen(2))
			Expect(cmds).To(ContainElement(cmd.Commands[0]))
			Expect(cmds).To(ContainElement(cmd.Commands[2]))
		})
	})

	Describe("VisibleCategories", func() {
		It("returns the visible categories", func() {
			categories := cmd.VisibleCategories()
			Expect(categories).To(HaveLen(2))

			cmds := categories[0].VisibleCommands()
			Expect(cmds).To(HaveLen(1))
			Expect(cmds).To(ContainElement(cmd.Commands[0]))

			cmds = categories[1].VisibleCommands()
			Expect(cmds).To(HaveLen(1))
			Expect(cmds).To(ContainElement(cmd.Commands[2]))
		})
	})
})

var _ = Describe("CommandsByName", func() {
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

var _ = Describe("CommandCategory", func() {
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
