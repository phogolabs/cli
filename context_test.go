package cli_test

import (
	"net"
	"net/url"
	"time"

	"github.com/phogolabs/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	var (
		context *cli.Context
		parent  *cli.Context
	)

	BeforeEach(func() {
		uri, err := url.Parse("http://example.com")
		Expect(err).To(BeNil())

		mac, err := net.ParseMAC("01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00")
		Expect(err).To(BeNil())

		parent = &cli.Context{
			Command: &cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level, l",
						Value: "debug",
					},
					&cli.StringSliceFlag{
						Name:  "user, u",
						Value: []string{"guest"},
					},
					&cli.BoolFlag{
						Name:  "verbose, v",
						Value: false,
					},
					&cli.IntFlag{
						Name:  "int-flag",
						Value: 2,
					},
					&cli.Int64Flag{
						Name:  "int64-flag",
						Value: 2,
					},
					&cli.UIntFlag{
						Name:  "uint-flag",
						Value: 2,
					},
					&cli.UInt64Flag{
						Name:  "uint64-flag",
						Value: 2,
					},
					&cli.Float32Flag{
						Name:  "float32-flag",
						Value: 2,
					},
					&cli.Float64Flag{
						Name:  "float64-flag",
						Value: 2,
					},
					&cli.URLFlag{
						Name:  "url-flag",
						Value: uri,
					},
					&cli.TimeFlag{
						Name:  "time-flag",
						Value: time.Now(),
					},
					&cli.DurationFlag{
						Name:  "duration-flag",
						Value: 20 * time.Second,
					},
					&cli.IPFlag{
						Name:  "ip-flag",
						Value: net.ParseIP("198.0.0.1"),
					},
					&cli.HardwareAddrFlag{
						Name:  "mac-flag",
						Value: mac,
					},
				},
			},
		}

		uri, err = url.Parse("http://google.com")
		Expect(err).To(BeNil())

		mac, err = net.ParseMAC("01-23-45-67-89-ab-cd-ef-00-00-01-23-45-67-89-ab-cd-ef-00-00")
		Expect(err).To(BeNil())

		context = &cli.Context{
			Parent: parent,
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
					&cli.IntFlag{
						Name:  "int-flag",
						Value: 1,
					},
					&cli.Int64Flag{
						Name:  "int64-flag",
						Value: 1,
					},
					&cli.UIntFlag{
						Name:  "uint-flag",
						Value: 1,
					},
					&cli.UInt64Flag{
						Name:  "uint64-flag",
						Value: 1,
					},
					&cli.Float32Flag{
						Name:  "float32-flag",
						Value: 1,
					},
					&cli.Float64Flag{
						Name:  "float64-flag",
						Value: 1,
					},
					&cli.URLFlag{
						Name:  "url-flag",
						Value: uri,
					},
					&cli.TimeFlag{
						Name:  "time-flag",
						Value: time.Now(),
					},
					&cli.DurationFlag{
						Name:  "duration-flag",
						Value: 10 * time.Second,
					},
					&cli.IPFlag{
						Name:  "ip-flag",
						Value: net.ParseIP("127.0.0.1"),
					},
					&cli.HardwareAddrFlag{
						Name:  "mac-flag",
						Value: mac,
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
				Expect(context.StringSlice("unknown")).To(BeEmpty())
			})
		})
	})

	Describe("StringSliceGlobal", func() {
		It("returns the value", func() {
			Expect(context.GlobalStringSlice("user")).To(ContainElement("guest"))
			Expect(context.GlobalStringSlice("u")).To(ContainElement("guest"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalStringSlice("unknown")).To(BeEmpty())
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

	Describe("Int", func() {
		It("returns the value", func() {
			Expect(context.Int("int-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Int("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalInt", func() {
		It("returns the value", func() {
			Expect(context.GlobalInt("int-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalInt("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("Int64", func() {
		It("returns the value", func() {
			Expect(context.Int64("int64-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Int64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalInt64", func() {
		It("returns the value", func() {
			Expect(context.GlobalInt64("int64-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalInt64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("UInt", func() {
		It("returns the value", func() {
			Expect(context.UInt("uint-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.UInt("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalUInt", func() {
		It("returns the value", func() {
			Expect(context.GlobalUInt("uint-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalUInt("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("UInt64", func() {
		It("returns the value", func() {
			Expect(context.UInt64("uint64-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.UInt64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalUInt64", func() {
		It("returns the value", func() {
			Expect(context.GlobalUInt64("uint64-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalUInt64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("Float32", func() {
		It("returns the value", func() {
			Expect(context.Float32("float32-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Float32("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalFloat32", func() {
		It("returns the value", func() {
			Expect(context.GlobalFloat32("float32-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalFloat32("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("Float64", func() {
		It("returns the value", func() {
			Expect(context.Float64("float64-flag")).To(BeNumerically("==", 1))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Float64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GlobalFloat64", func() {
		It("returns the value", func() {
			Expect(context.GlobalFloat64("float64-flag")).To(BeNumerically("==", 2))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalFloat64("unknown")).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("URL", func() {
		It("returns the value", func() {
			Expect(context.URL("url-flag").String()).To(Equal("http://google.com"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.URL("unknown")).To(BeNil())
			})
		})
	})

	Describe("GlobalURL", func() {
		It("returns the value", func() {
			Expect(context.GlobalURL("url-flag").String()).To(Equal("http://example.com"))
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalURL("unknown")).To(BeNil())
			})
		})
	})

	Describe("Get", func() {
		It("returns the value", func() {
			Expect(context.Get("url-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Get("unknown")).To(BeNil())
			})
		})
	})

	Describe("GlobalGet", func() {
		It("returns the value", func() {
			Expect(context.GlobalGet("url-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalGet("unknown")).To(BeNil())
			})
		})
	})

	Describe("Time", func() {
		It("returns the value", func() {
			Expect(context.Time("time-flag")).NotTo(BeZero())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Time("unknown")).To(BeZero())
			})
		})
	})

	Describe("GlobalTime", func() {
		It("returns the value", func() {
			Expect(context.GlobalTime("time-flag")).NotTo(BeZero())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalTime("unknown")).To(BeZero())
			})
		})
	})

	Describe("Duration", func() {
		It("returns the value", func() {
			Expect(context.Duration("duration-flag")).NotTo(BeZero())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.Duration("unknown")).To(BeZero())
			})
		})
	})

	Describe("GlobalDuration", func() {
		It("returns the value", func() {
			Expect(context.GlobalDuration("duration-flag")).NotTo(BeZero())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalDuration("unknown")).To(BeZero())
			})
		})
	})

	Describe("IP", func() {
		It("returns the value", func() {
			Expect(context.IP("ip-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.IP("unknown")).To(BeNil())
			})
		})
	})

	Describe("GlobalIP", func() {
		It("returns the value", func() {
			Expect(context.GlobalIP("ip-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalIP("unknown")).To(BeZero())
			})
		})
	})

	Describe("HardwareAddr", func() {
		It("returns the value", func() {
			Expect(context.HardwareAddr("mac-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.HardwareAddr("unknown")).To(BeNil())
			})
		})
	})

	Describe("GlobalHardwareAddr", func() {
		It("returns the value", func() {
			Expect(context.GlobalHardwareAddr("mac-flag")).NotTo(BeNil())
		})

		Context("when the flag cannot be found", func() {
			It("returns default value", func() {
				Expect(context.GlobalHardwareAddr("unknown")).To(BeZero())
			})
		})
	})
})

// &cli.HardwareAddrFlag
