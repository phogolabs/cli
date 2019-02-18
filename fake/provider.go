// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"sync"

	"github.com/phogolabs/cli"
)

type Provider struct {
	ProvideStub        func(*cli.Context) error
	provideMutex       sync.RWMutex
	provideArgsForCall []struct {
		arg1 *cli.Context
	}
	provideReturns struct {
		result1 error
	}
	provideReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Provider) Provide(arg1 *cli.Context) error {
	fake.provideMutex.Lock()
	ret, specificReturn := fake.provideReturnsOnCall[len(fake.provideArgsForCall)]
	fake.provideArgsForCall = append(fake.provideArgsForCall, struct {
		arg1 *cli.Context
	}{arg1})
	fake.recordInvocation("Provide", []interface{}{arg1})
	fake.provideMutex.Unlock()
	if fake.ProvideStub != nil {
		return fake.ProvideStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.provideReturns.result1
}

func (fake *Provider) ProvideCallCount() int {
	fake.provideMutex.RLock()
	defer fake.provideMutex.RUnlock()
	return len(fake.provideArgsForCall)
}

func (fake *Provider) ProvideArgsForCall(i int) *cli.Context {
	fake.provideMutex.RLock()
	defer fake.provideMutex.RUnlock()
	return fake.provideArgsForCall[i].arg1
}

func (fake *Provider) ProvideReturns(result1 error) {
	fake.ProvideStub = nil
	fake.provideReturns = struct {
		result1 error
	}{result1}
}

func (fake *Provider) ProvideReturnsOnCall(i int, result1 error) {
	fake.ProvideStub = nil
	if fake.provideReturnsOnCall == nil {
		fake.provideReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.provideReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Provider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.provideMutex.RLock()
	defer fake.provideMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Provider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ cli.Provider = new(Provider)
