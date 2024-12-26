package libs

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrVerificationFailed     = errors.New("token verification failed")
	PIPEOPS_CONTROL_PLANE_API = ""
)

type HttpClients interface {
	VerifyToken(token string, operatorID string) (*models.PipeOpsTokenVerificationResponse, error)
}

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient() HttpClients {
	r := resty.New()
	r.Debug = true

	URL := strings.TrimSpace(PIPEOPS_CONTROL_PLANE_API)
	r.SetBaseURL(URL)

	return &HttpClient{
		client: r,
	}
}

// VerifyToken performs a POST request to verify a token.
func (v *HttpClient) VerifyToken(token string, operatorID string) (*models.PipeOpsTokenVerificationResponse, error) {
	if strings.TrimSpace(token) == "" || strings.TrimSpace(operatorID) == "" {
		return nil, errors.New("token or operatorID is empty")
	}

	payload := map[string]string{
		"token":       token,
		"operator_id": operatorID,
	}

	resp, err := v.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 401 || resp.StatusCode() == 400 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, ErrVerificationFailed
	}

	var respData *models.PipeOpsTokenVerificationResponse
	if err := json.Unmarshal(resp.Body(), &respData); err != nil {
		return nil, err
	}

	if !respData.Valid {
		return nil, ErrInvalidToken
	}

	return respData, nil
}
