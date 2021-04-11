package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/controller"
	"github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/event"
	"gopkg.in/yaml.v2"
)

type Handler struct {
	ctrl *controller.Controller
}

const (
	systemUserCreatedMsg = "AWS Account for system user has been created: `{{.UserName}}`"
	userCreatedMsg       = `Your AWS Account has been created.
Initial password:

` + "```" + `
{{.Password}}
` + "```" + `

Please sign in AWS and change your password.

https://signin.aws.amazon.com/console

Please create your AWS Access Key if needed.
https://console.aws.amazon.com/iam/home#/users/{{.UserName}}?section=security_credentials
`
)

func (handler *Handler) Init(ctx context.Context) error {
	cfgString := os.Getenv("CONFIG")
	cfg := controller.Config{}
	if cfgString != "" {
		if err := yaml.Unmarshal([]byte(cfgString), &cfg); err != nil {
			return fmt.Errorf("unmarshal config: %w", err)
		}
	}

	// set default config
	if cfg.InitialPasswordLength == 0 {
		cfg.InitialPasswordLength = 32
	}
	if cfg.Message == "" {
		cfg.Message = userCreatedMsg
	}
	if cfg.MessageForSystemUser == "" {
		cfg.MessageForSystemUser = systemUserCreatedMsg
	}

	sess := session.Must(session.NewSession())

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(cfg.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(cfg.SecretID),
	}
	if cfg.VersionID != "" {
		input.VersionId = aws.String(cfg.VersionID)
	}
	output, err := svc.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("get secret value from AWS SecretsManager: %w", err)
	}
	secret := controller.Secret{}
	if err := yaml.Unmarshal([]byte(*output.SecretString), &secret); err != nil {
		return fmt.Errorf("parse secret value: %w", err)
	}

	ctrl, _, err := controller.New(ctx, controller.Param{})
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}
	ctrl.Config = cfg

	if ctrl.MessageTemplate, err = ctrl.CompileTemplate(cfg.Message); err != nil {
		return err
	}
	if ctrl.MessageTemplateForSystemUser, err = ctrl.CompileTemplate(cfg.MessageForSystemUser); err != nil {
		return err
	}

	// create a slack client
	ctrl.SlackBot = slack.New(secret.SlackBotToken)
	handler.ctrl = &ctrl

	return nil
}

func (handler *Handler) Start(ctx context.Context, ev events.CloudWatchEvent) error {
	if err := json.NewEncoder(os.Stderr).Encode(&ev); err != nil {
		return fmt.Errorf("parse a CloudWatchEvent as JSON: %w", err)
	}
	if err := handler.start(ctx, ev); err != nil {
		logrus.WithError(err).Error("start")
		return err
	}
	return nil
}

func (handler *Handler) start(ctx context.Context, ev events.CloudWatchEvent) error {
	parser := event.Parser{}
	user, err := parser.Parse(ctx, ev)
	if err != nil {
		return fmt.Errorf("parse a CloudWatchEvent: %w", err)
	}

	param := controller.Param{
		UserName: user.Name,
	}

	return handler.ctrl.Run(ctx, param)
}
