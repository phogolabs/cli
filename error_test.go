package cli_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/cli"
)

var _ = Describe("ExitErrorCollector", func() {
	It("creates a new error", func() {
		err := cli.ExitErrorCollector{
			fmt.Errorf("oh no!"),
		}

		Expect(err).To(HaveLen(1))
		Expect(err[0]).To(MatchError("oh no!"))
	})

	It("formats the error", func() {
		err := cli.ExitErrorCollector{
			fmt.Errorf("oh no!"),
			fmt.Errorf("oh ye!"),
		}

		Expect(err).To(MatchError("oh no!\noh ye!"))
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
		errx := cli.WrapError(err, 1)
		Expect(errx.Unwrap()).To(Equal(err))
	})
})
