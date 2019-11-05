package interactors

import (
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
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
		email      string
		password   string
		fakeLogger *FakeCreateUserLogger
		expected   string
	}{
		{
			"shouldCallLogger",
			"fake@email.com",
			"123",
			&FakeCreateUserLogger{},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createUserInteractor := NewCreateUser(tc.fakeLogger)
			user := entities.User{Email: tc.email, Password: tc.password}
			err := createUserInteractor.Execute(user)
			if err != nil && err.Error() != tc.expected {
				t.Errorf("Create User failed. Error: %s", err)
				return
			}
			t.Logf("Create user ok")
		})
	}

}
