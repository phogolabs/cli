package s3_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws/s3"
	"github.com/phogolabs/cli/provider/aws/s3/fake"
)

var _ = Describe("Provider", func() {
	var (
		flag *cli.StringFlag
		ctx  *cli.Context
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:     "listen-addr",
			Usage:    "listen address of HTTP server",
			FilePath: "s3://*.txt",
		}

		ctx = &cli.Context{
			Command: &cli.Command{
				Name:  "app",
				Flags: []cli.Flag{flag},
			},
		}
	})

	Describe("S3", func() {
		var (
			provider   *s3.Provider
			fileSystem *fake.FileSystem
		)

		BeforeEach(func() {
			fileSystem = &fake.FileSystem{}
			fileSystem.GlobReturns([]string{"report.txt"}, nil)
			fileSystem.ReadFileReturns([]byte("9292"), nil)

			provider = &s3.Provider{
				FileSystem: fileSystem,
			}
		})

		It("sets the value successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("9292"))

			Expect(fileSystem.GlobCallCount()).To(Equal(1))
			Expect(fileSystem.GlobArgsForCall(0)).To(Equal("*.txt"))

			Expect(fileSystem.ReadFileCallCount()).To(Equal(1))
			Expect(fileSystem.ReadFileArgsForCall(0)).To(Equal("report.txt"))
		})

		Context("when the globbing fails", func() {
			BeforeEach(func() {
				fileSystem.GlobReturns(nil, fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(provider.Provide(ctx)).To(MatchError("oh no!"))
			})
		})

		Context("when the file path is not s3", func() {
			BeforeEach(func() {
				flag.FilePath = "*.docx"
			})

			It("does not set the value", func() {
				Expect(provider.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})

		Context("when the file path is not set", func() {
			BeforeEach(func() {
				flag.FilePath = ""
			})

			It("does not set the value", func() {
				Expect(provider.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})
	})
})
