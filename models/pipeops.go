package models

// TokenVerificationResponse represents the response structure from PipeOps API
type PipeOpsTokenVerificationResponse struct {
	Valid    bool   `json:"valid"`
	NATSURL  string `json:"nats_url"`
	ErrorMsg string `json:"error_msg,omitempty"`
}
