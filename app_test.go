package cli_test

import (
	"bytes"
	"os"

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

		flags := []cli.Flag{
			&cli.StringFlag{
				Name: "test.timeout",
			},
			&cli.StringFlag{
				Name: "test.coverprofile",
			},
			&cli.StringFlag{
				Name: "ginkgo.seed",
			},
			&cli.StringFlag{
				Name: "ginkgo.slowSpecThreshold",
			},
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
			Flags:    flags,
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

			Expect(cmd.Metadata).To(HaveKeyWithValue("HideVersion", app.HideVersion))
			Expect(cmd.Metadata).To(HaveKeyWithValue("Version", app.Version))
			Expect(cmd.Metadata).To(HaveKeyWithValue("Authors", app.Authors))
			Expect(cmd.Metadata).To(HaveKeyWithValue("Copyright", app.Copyright))

			Expect(ctx.Args).To(BeEmpty())

			return nil
		}

		Expect(app.Run(os.Args)).To(Succeed())
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
