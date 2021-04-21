package lambda

import (
	"errors"

	"github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/controller"
)

var (
	errAWSAccountIDIsRequired = errors.New("AWS Account ID is required")
	errSecretIDIsRequired     = errors.New("secret ID is required")
)

func (handler *Handler) validateConfig(cfg controller.Config) error {
	if cfg.AWSAccountID == "" {
		return errAWSAccountIDIsRequired
	}
	if cfg.SecretID == "" {
		return errSecretIDIsRequired
	}
	return nil
}
