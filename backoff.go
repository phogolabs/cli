package cli

import (
	"time"

	"github.com/cenkalti/backoff"
	"github.com/phogolabs/log"
)

// BackOffStrategy represents the backoff strategy
type BackOffStrategy backoff.BackOff

// NewExponentialBackOffStrategy creates a new backoff strategy
func NewExponentialBackOffStrategy() BackOffStrategy {
	strategy := backoff.NewExponentialBackOff()
	strategy.MaxElapsedTime = 30 * time.Second
	strategy.InitialInterval = 2 * time.Second
	return strategy
}

// BackOffProvider backoff the provider
type BackOffProvider struct {
	Provider Provider
	Strategy BackOffStrategy
}

// Provide parses the args
func (m *BackOffProvider) Provide(ctx *Context) error {
	tryProvide := func() error {
		log.Info("providing the application arguments")
		return m.Provider.Provide(ctx)
	}

	notify := func(err error, t time.Duration) {
		log.WithError(err).Warnf("providing the application arguments not successful. retry in %v", t)
	}

	if m.Strategy == nil {
		// set the default strategy
		m.Strategy = NewExponentialBackOffStrategy()
	}

	if err := backoff.RetryNotify(tryProvide, m.Strategy, notify); err != nil {
		log.WithError(err).Fatal("providing the application argument failed")
		return err
	}

	return nil
}
