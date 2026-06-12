package ds

import (
	"slices"

	model "tictactoe/internal/domain/model"
)

type SessionConfig struct {
	UUIDs []string
	State *string
}

type SessionOpt interface {
	ApplyToSession(*SessionConfig)
}

type UserConfig struct {
	UUIDs  []string
	Logins []string
}

type UserOpt interface {
	ApplyToUser(*UserConfig)
}

type withUUID struct{ uuids []string }

func (o withUUID) ApplyToSession(c *SessionConfig) { c.UUIDs = o.uuids }
func (o withUUID) ApplyToUser(c *UserConfig)       { c.UUIDs = o.uuids }

func WithUUID(UUIDs ...string) withUUID { return withUUID{uuids: slices.Clone(UUIDs)} }

type withState struct{ state string }

func (o withState) ApplyToSession(c *SessionConfig) {
	c.State = &o.state
}

func WithState(state model.State) withState {
	return withState{state: state.String()}
}

type withLogin struct{ logins []string }

func (o withLogin) ApplyToUser(c *UserConfig) { c.Logins = o.logins }

func WithLogin(logins ...string) withLogin {
	return withLogin{logins: slices.Clone(logins)}
}
