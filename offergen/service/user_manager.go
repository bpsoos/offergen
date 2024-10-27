package service

import (
	"offergen/logging"
)

type Manager struct {
	persister     UserPersister
	authenticator Authenticator
}

type UserPersister interface {
	Save(id, address string) error
	GetEmail(id string) (string, error)
	Delete(id string) error
}

type Authenticator interface {
	DeleteUser(authToken []byte) error
}

type UserManagerDeps struct {
	Persister     UserPersister
	Authenticator Authenticator
}

func NewUserManager(deps *UserManagerDeps) *Manager {
	return &Manager{
		persister:     deps.Persister,
		authenticator: deps.Authenticator,
	}
}

var logger = logging.GetLogger()

func (m *Manager) Save(id, email string) error {
	logger.Info("saving user", "userID", id)

	return m.persister.Save(id, email)
}

func (m *Manager) GetEmail(id string) (string, error) {
	return m.persister.GetEmail(id)
}

func (m *Manager) Delete(id string, token []byte) error {
	if err := m.authenticator.DeleteUser(token); err != nil {
		return err
	}

	return m.persister.Delete(id)
}
