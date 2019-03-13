package cli_test

import (
	"bytes"
	"fmt"
	"os"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
)

var _ = Describe("App", func() {
	var app *cli.App

	BeforeEach(func() {

		noop := func(ctx *cli.Context) error {
			return nil
		}

		commands := []*cli.Command{
			&cli.Command{
				Name:        "sync",
				Usage:       "Generate a SQL script of CRUD operations for given database schema",
				Description: "Generate a SQL script of CRUD operations for given database schema",
			},
		}

		app = &cli.App{
			Name:        "prana",
			Usage:       "Golang Database Manager",
			UsageText:   "prana [global options]",
			HelpName:    "prana",
			Description: "Golang Database Manager",
			ArgsUsage:   "[args]",
			Version:     "1.0-beta-04",
			Copyright:   "Open Source",
			Authors: []*cli.Author{
				&cli.Author{
					Name:  "John Freeman",
					Email: "john@exmaple.com",
				},
			},
			Before:   noop,
			After:    noop,
			Action:   nil,
			Commands: commands,
		}
	})

	It("executes the app's command successfully", func() {
		app.Action = func(ctx *cli.Context) error {
			cmd := ctx.Command

			Expect(cmd).NotTo(BeNil())

			Expect(cmd.Name).To(Equal(app.Name))
			Expect(cmd.Usage).To(Equal(app.Usage))
			Expect(cmd.UsageText).To(Equal(app.UsageText))
			Expect(cmd.HideHelp).To(Equal(app.HideHelp))
			Expect(cmd.HelpName).To(Equal(app.HelpName))
			Expect(cmd.Description).To(Equal(app.Description))
			Expect(cmd.ArgsUsage).To(Equal(app.ArgsUsage))

			Expect(cmd.Commands).To(ContainElement(cli.HelpCommand))

			Expect(ctx.Manifest.HideVersion).To(Equal(app.HideVersion))
			Expect(ctx.Manifest.Version).To(Equal(app.Version))
			Expect(ctx.Manifest.Copyright).To(Equal(app.Copyright))
			Expect(ctx.Manifest.Authors).To(Equal(app.Authors))

			Expect(ctx.Args).To(BeEmpty())

			return nil
		}

		Expect(app.Run([]string{"app"})).To(Succeed())
	})

	Context("when the operation system sends a signal", func() {
		It("handles the signal", func() {
			count := 0

			app.Signals = []os.Signal{syscall.SIGUSR1}
			app.Action = func(ctx *cli.Context) error {
				return nil
			}

			app.OnSignal = func(ctx *cli.Context) error {
				count++

				cmd := ctx.Command
				Expect(cmd).NotTo(BeNil())
				Expect(cmd.Name).To(Equal(app.Name))

				Expect(ctx.Signal).To(Equal(syscall.SIGUSR1))
				return nil
			}

			Expect(app.Run([]string{"app"})).To(Succeed())

			process, err := os.FindProcess(os.Getpid())
			Expect(err).NotTo(HaveOccurred())
			Expect(process.Signal(syscall.SIGUSR1)).To(Succeed())

			Eventually(func() int { return count }).Should(Equal(1))
		})
	})

	Context("when the app name is not provided", func() {
		It("sets the app name", func() {
			app.Name = ""

			app.Action = func(ctx *cli.Context) error {
				cmd := ctx.Command
				Expect(cmd.Name).To(Equal(app.Name))
				Expect(cmd.Name).To(Equal("app"))
				return nil
			}

			Expect(app.Run([]string{"app"})).To(Succeed())
		})
	})

	Context("when the app version is requested", func() {
		var buffer *bytes.Buffer

		BeforeEach(func() {
			buffer = &bytes.Buffer{}
			app.Writer = buffer
		})

		It("shows the version", func() {
			Expect(app.Run([]string{"app", "-v"})).To(Succeed())
			Expect(buffer.String()).To(Equal("prana version 1.0-beta-04\n"))
		})
	})

	Context("when the app fails", func() {
		It("exits with the provided code", func() {
			app.Action = func(ctx *cli.Context) error {
				return cli.NewExitError("oh no!", 78)
			}

			app.Exit = func(code int) {
				Expect(code).To(Equal(78))
			}

			Expect(app.Run([]string{"app"})).To(MatchError("oh no!"))
		})

		Context("when the error is not exit error", func() {
			It("exits with the exit code 1", func() {
				app.Action = func(ctx *cli.Context) error {
					return fmt.Errorf("oh no!")
				}

				app.OnExitErr = func(err error) error {
					Expect(err).To(MatchError("oh no!"))
					return err
				}

				app.Exit = func(code int) {
					Expect(code).To(Equal(1))
				}

				Expect(app.Run([]string{"app"})).To(MatchError("oh no!"))
			})
		})
	})
})

var _ = Describe("Author", func() {
	It("returns the author as string", func() {
		author := &cli.Author{
			Name: "John",
		}

		Expect(author.String()).To(Equal("John"))
	})

	Context("when the author has email", func() {
		It("returns the author as string", func() {
			author := &cli.Author{
				Name:  "John",
				Email: "john@example.com",
			}

			Expect(author.String()).To(Equal("John <john@example.com>"))
		})
	})
})
