package cli_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/cli"
)

var _ = Describe("MultiError", func() {
	It("creates a new error", func() {
		m := cli.NewMultiError(fmt.Errorf("oh no!"))
		Expect(*m).To(HaveLen(1))
		Expect((*m)[0]).To(MatchError("oh no!"))
	})

	It("returns the exit code", func() {
		m := cli.NewMultiError(fmt.Errorf("oh no!"))
		Expect(m.ExitCode()).To(Equal(1))
	})

	Context("when the underlying error is exit error", func() {
		It("returns the exit code", func() {
			m := cli.NewMultiError(cli.NewExitError("oh no", 88))
			Expect(m.ExitCode()).To(Equal(88))
		})
	})

	It("formats the error", func() {
		err := &cli.MultiError{
			fmt.Errorf("oh no!"),
			fmt.Errorf("oh ye!"),
		}

		Expect(err).To(MatchError("oh no!\noh ye!"))
	})

	Describe("AppendError", func() {
		It("appends error to a multi error", func() {
			m := cli.NewMultiError(fmt.Errorf("oh no!"))
			err := cli.AppendError(m, fmt.Errorf("oh ye!"))
			Expect(err).To(MatchError("oh no!\noh ye!"))
		})

		Context("when the target error is not multi error", func() {
			It("creates a multi error", func() {
				err := cli.AppendError(fmt.Errorf("oh no!"), fmt.Errorf("oh ye!"))
				Expect(err).To(MatchError("oh no!\noh ye!"))
			})
		})

		Context("when the target error is nil", func() {
			It("returns the error", func() {
				err := fmt.Errorf("oh no!")
				Expect(cli.AppendError(nil, err)).To(Equal(err))
			})
		})
	})
})

var _ = Describe("ExitError", func() {
	It("creates a new exit error", func() {
		err := cli.NewExitError("oh no!", 69)
		Expect(err.Error()).To(Equal("oh no!"))
		Expect(err.ExitCode()).To(Equal(69))
	})

	It("wraps an error", func() {
		err := cli.WrapExitError(fmt.Errorf("oh no!"), 69)
		Expect(err.Error()).To(Equal("oh no!"))
		Expect(err.ExitCode()).To(Equal(69))
	})
})
