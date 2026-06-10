package dto

type SessionResponse struct {
	SessionUUID string    `json:"session_uuid"`
	State       string    `json:"state"`
	Board       []string  `json:"board"`
	Winner      string    `json:"winner,omitempty"`
	Players     []*Player `json:"players"`
}

type Player struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Mark string `json:"mark"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	UUID string `json:"uuid"`
}
