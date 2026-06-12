package port

import (
	"slices"

	"tictactoe/internal/domain/model"
)

type SessionGetConfig struct {
	UUIDs []string
	State *model.State
}
type SessionGetOpt interface{ ApplyToSession(*SessionGetConfig) }

type RulesGetConfig struct {
	UUIDs []string
}
type RulesGetOpt interface{ ApplyToRules(*RulesGetConfig) }

type UserGetConfig struct {
	UUIDs  []string
	Logins []string
}
type UserGetOpt interface{ ApplyToUser(*UserGetConfig) }

type withUUID struct{ uuids []string }

func (o withUUID) ApplyToSession(c *SessionGetConfig) { c.UUIDs = o.uuids }
func (o withUUID) ApplyToRules(c *RulesGetConfig)     { c.UUIDs = o.uuids }
func (o withUUID) ApplyToUser(c *UserGetConfig)       { c.UUIDs = o.uuids }

func WithUUID(UUIDs ...string) withUUID {
	return withUUID{uuids: slices.Clone(UUIDs)}
}

type withState struct{ state model.State }

func WithState(s model.State) withState {
	return withState{state: s}
}

func (o withState) ApplyToSession(c *SessionGetConfig) { c.State = &o.state }

type withLogin struct{ logins []string }

func (o withLogin) ApplyToUser(c *UserGetConfig) { c.Logins = o.logins }

func WithLogin(logins ...string) withLogin {
	return withLogin{logins: slices.Clone(logins)}
}
