package interactors

import (
	"testing"
)

type FakeCreateUserLogger struct {
}

func (fl *FakeCreateUserLogger) Info(...interface{}) {}

func (fl *FakeCreateUserLogger) Infof(string, ...interface{}) {}

func (fl *FakeCreateUserLogger) Debug(...interface{}) {}

func (fl *FakeCreateUserLogger) Warn(...interface{}) {}

func (fl *FakeCreateUserLogger) Error(...interface{}) {}

func (fl *FakeCreateUserLogger) Errorf(string, ...interface{}) {}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name       string
		fakeLogger *FakeCreateUserLogger
		expected   string
	}{
		{
			"shouldCallLogger",
			&FakeCreateUserLogger{},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createUserInteractor := NewCreateUser(tc.fakeLogger)
			createUserInteractor.Execute()
			t.Logf("Create user ok")
		})
	}

}
