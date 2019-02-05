package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/cli"
)

var _ = Describe("Context", func() {
	var context *cli.Context

	BeforeEach(func() {
		context = &cli.Context{
			Command: &cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level, l",
						Value: "info",
					},
					&cli.StringSliceFlag{
						Name:  "user, u",
						Value: []string{"root"},
					},
					&cli.BoolFlag{
						Name:  "verbose, v",
						Value: true,
					},
				},
			},
		}
	})

	Describe("String", func() {
		It("returns the value", func() {
			Expect(context.String("log-level")).To(Equal("info"))
			Expect(context.String("l")).To(Equal("info"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.String("unknown")).To(BeEmpty())
			})
		})
	})

	Describe("GloablString", func() {
		BeforeEach(func() {
			context.Parent = &cli.Context{
				Command: &cli.Command{
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "log-level, l",
							Value: "debug",
						},
					},
				},
			}
		})

		It("returns the value", func() {
			Expect(context.GlobalString("log-level")).To(Equal("debug"))
			Expect(context.GlobalString("l")).To(Equal("debug"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalString("unknown")).To(BeEmpty())
			})
		})
	})

	Describe("StringSlice", func() {
		It("returns the value", func() {
			Expect(context.StringSlice("user")).To(ContainElement("root"))
			Expect(context.StringSlice("u")).To(ContainElement("root"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.StringSlice("unknonw")).To(BeEmpty())
			})
		})
	})

	Describe("StringSliceGlobal", func() {
		BeforeEach(func() {
			context.Parent = &cli.Context{
				Command: &cli.Command{
					Flags: []cli.Flag{
						&cli.StringSliceFlag{
							Name:  "user, u",
							Value: []string{"guest"},
						},
					},
				},
			}
		})

		It("returns the value", func() {
			Expect(context.GlobalStringSlice("user")).To(ContainElement("guest"))
			Expect(context.GlobalStringSlice("u")).To(ContainElement("guest"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalStringSlice("unknonw")).To(BeEmpty())
			})
		})
	})

	Describe("Bool", func() {
		It("returns the value", func() {
			Expect(context.Bool("verbose")).To(BeTrue())
			Expect(context.Bool("v")).To(BeTrue())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Bool("unknown")).To(BeFalse())
			})
		})
	})

	Describe("GlobalBool", func() {
		BeforeEach(func() {
			context.Parent = &cli.Context{
				Command: &cli.Command{
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "verbose, v",
							Value: false,
						},
					},
				},
			}
		})

		It("returns the value", func() {
			Expect(context.GlobalBool("verbose")).To(BeFalse())
			Expect(context.GlobalBool("v")).To(BeFalse())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalBool("unknown")).To(BeFalse())
			})
		})
	})
})
