// Code generated by counterfeiter. DO NOT EDIT.
package repofakes

import (
	"sync"

	"github.com/devrajsinghrawat/invite_app/src/model"
	"github.com/devrajsinghrawat/invite_app/src/repo"
)

type FakeIUser struct {
	GetUserStub        func(*model.User) (model.User, error)
	getUserMutex       sync.RWMutex
	getUserArgsForCall []struct {
		arg1 *model.User
	}
	getUserReturns struct {
		result1 model.User
		result2 error
	}
	getUserReturnsOnCall map[int]struct {
		result1 model.User
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeIUser) GetUser(arg1 *model.User) (model.User, error) {
	fake.getUserMutex.Lock()
	ret, specificReturn := fake.getUserReturnsOnCall[len(fake.getUserArgsForCall)]
	fake.getUserArgsForCall = append(fake.getUserArgsForCall, struct {
		arg1 *model.User
	}{arg1})
	stub := fake.GetUserStub
	fakeReturns := fake.getUserReturns
	fake.recordInvocation("GetUser", []interface{}{arg1})
	fake.getUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeIUser) GetUserCallCount() int {
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	return len(fake.getUserArgsForCall)
}

func (fake *FakeIUser) GetUserCalls(stub func(*model.User) (model.User, error)) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = stub
}

func (fake *FakeIUser) GetUserArgsForCall(i int) *model.User {
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	argsForCall := fake.getUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeIUser) GetUserReturns(result1 model.User, result2 error) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = nil
	fake.getUserReturns = struct {
		result1 model.User
		result2 error
	}{result1, result2}
}

func (fake *FakeIUser) GetUserReturnsOnCall(i int, result1 model.User, result2 error) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = nil
	if fake.getUserReturnsOnCall == nil {
		fake.getUserReturnsOnCall = make(map[int]struct {
			result1 model.User
			result2 error
		})
	}
	fake.getUserReturnsOnCall[i] = struct {
		result1 model.User
		result2 error
	}{result1, result2}
}

func (fake *FakeIUser) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeIUser) recordInvocation(key string, args []interface{}) {
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

var _ repo.IUser = new(FakeIUser)
