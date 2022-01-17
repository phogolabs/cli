package cli_test

import (
	"fmt"

	"github.com/phogolabs/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExitErrorCollector", func() {
	It("creates a new error", func() {
		err := cli.ExitErrorCollector{
			fmt.Errorf("oh no"),
		}

		Expect(err).To(HaveLen(1))
		Expect(err[0]).To(MatchError("oh no"))
	})

	It("formats the error", func() {
		err := cli.ExitErrorCollector{
			fmt.Errorf("oh no"),
			fmt.Errorf("oh ye"),
		}

		Expect(err).To(MatchError("oh no\noh ye"))
	})
})

var _ = Describe("ExitError", func() {
	It("creates a new exit error", func() {
		err := cli.NewExitError("oh no!", 69)
		Expect(err.Error()).To(Equal("oh no!"))
		Expect(err.Code()).To(Equal(69))
	})

	It("wraps an error", func() {
		err := fmt.Errorf("oh no")
		errx := cli.WrapError(err)
		Expect(errx.Unwrap()).To(Equal(err))
	})

	Context("WithCode", func() {
		It("returns a error with new code", func() {
			err := cli.NewExitError("oh no", 69)
			errx := err.WithCode(129)

			Expect(errx).To(MatchError("oh no"))
			Expect(errx.Code()).To(Equal(129))
			Expect(errx.Code()).NotTo(Equal(err.Code()))
		})
	})
})
