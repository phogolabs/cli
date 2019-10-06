package aws_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws"
	"github.com/phogolabs/cli/provider/aws/fake"
)

var _ = Describe("S3Provider", func() {
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

		bucket := &cli.StringFlag{
			Name:  "aws-bucket",
			Value: "my-bucket",
		}

		ctx = &cli.Context{
			Command: &cli.Command{
				Name:  "app",
				Flags: []cli.Flag{flag, bucket},
			},
		}
	})

	Describe("S3", func() {
		var (
			provider   *aws.S3Provider
			fileSystem *fake.FileSystem
		)

		BeforeEach(func() {
			fileSystem = &fake.FileSystem{}
			fileSystem.GlobReturns([]string{"report.txt"}, nil)
			fileSystem.ReadFileReturns([]byte("9292"), nil)

			provider = &aws.S3Provider{
				FileSystem: fileSystem,
			}
		})

		It("sets the value successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("9292"))

			Expect(fileSystem.GlobCallCount()).To(Equal(1))

			bucket, pattern := fileSystem.GlobArgsForCall(0)
			Expect(bucket).To(Equal("my-bucket"))
			Expect(pattern).To(Equal("*.txt"))

			Expect(fileSystem.ReadFileCallCount()).To(Equal(1))

			bucket, file := fileSystem.ReadFileArgsForCall(0)
			Expect(bucket).To(Equal("my-bucket"))
			Expect(file).To(Equal("report.txt"))
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

var _ = Describe("SSMProvider", func() {
	var (
		flag *cli.StringFlag
		ctx  *cli.Context
	)

	BeforeEach(func() {
		flag = &cli.StringFlag{
			Name:  "listen-addr",
			Usage: "listen address of HTTP server",
			Metadata: cli.Map{
				"ssm_param": "/terraform/secret",
			},
		}

		ctx = &cli.Context{
			Command: &cli.Command{
				Name:  "app",
				Flags: []cli.Flag{flag},
			},
		}
	})

	Describe("SSM", func() {
		var (
			provider *aws.SSMProvider
			client   *fake.Client
		)

		BeforeEach(func() {
			client = &fake.Client{}
			client.GetReturns("swordfish", nil)

			provider = &aws.SSMProvider{
				Client: client,
			}
		})

		It("sets the value successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(flag.Value).To(Equal("swordfish"))

			Expect(client.GetCallCount()).To(Equal(1))
			Expect(client.GetArgsForCall(0)).To(Equal("/terraform/secret"))
		})

		Context("when the file path is not set", func() {
			BeforeEach(func() {
				flag.Metadata = cli.Map{}
			})

			It("does not set the value", func() {
				Expect(provider.Provide(ctx)).To(Succeed())
				Expect(flag.Value).To(BeZero())
			})
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.GetReturns("", fmt.Errorf("oh no!"))
			})

			It("returns an error", func() {
				Expect(provider.Provide(ctx)).To(MatchError("oh no!"))
			})
		})
	})
})
