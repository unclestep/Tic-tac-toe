package dto

type SessionResponse struct {
	SessionID string   `json:"session_id"`
	State     string   `json:"state"`
	Board     []int8   `json:"board"`
	Winner    string   `json:"winner,omitempty"`
	Players   []string `json:"players"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
