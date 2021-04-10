package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sethvargo/go-password/password"
	"github.com/slack-go/slack"
)

type Secret struct {
	SlackBotToken string `yaml:"slack_bot_token"`
}

func (ctrl *Controller) GetSlackUser(ctx context.Context, name string) (slack.User, error) {
	users, err := ctrl.SlackBot.GetUsersContext(ctx)
	if err != nil {
		return slack.User{}, fmt.Errorf("get all Slack Users: %w", err)
	}
	for _, user := range users {
		if user.Name == name {
			return user, nil
		}
	}
	return slack.User{}, errors.New("user isn't found: " + name)
}

func (ctrl *Controller) Run(ctx context.Context, param Param) error {
	sess := session.Must(session.NewSession())

	// get a slack user id
	user, err := ctrl.GetSlackUser(ctx, param.UserName)
	if err != nil {
		// treat the user as a system account
		// send a notification to slack
		if _, _, _, err := ctrl.SlackBot.SendMessageContext(ctx, ctrl.Config.Slack.ChannelIDForSystemUser, slack.MsgOptionText(ctrl.Config.MessageForSystemUser+param.UserName, false)); err != nil {
			return fmt.Errorf("send a notification that a system user has been created to Slack channel (channel id: %s): %w", ctrl.Config.Slack.ChannelIDForSystemUser, err)
		}
		return nil
	}
	// create an initial password
	passwd, err := password.Generate(ctrl.Config.InitialPasswordLength, 10, 10, false, false)
	if err != nil {
		return fmt.Errorf("general an initial password: %w", err)
	}
	// create a login profile
	iamSvc := iam.New(sess, aws.NewConfig().WithRegion(ctrl.Config.Region))
	if _, err := iamSvc.CreateLoginProfileWithContext(ctx, &iam.CreateLoginProfileInput{
		Password:              aws.String(passwd),
		PasswordResetRequired: aws.Bool(true),
		UserName:              aws.String(param.UserName),
	}); err != nil {
		return fmt.Errorf("create a login profile: %w", err)
	}
	// create an access key
	// if _, err := iamSvc.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{}); err != nil {
	// 	return err
	// }
	// create a message
	// send a message
	if _, _, _, err := ctrl.SlackBot.SendMessageContext(ctx, user.ID, slack.MsgOptionText(ctrl.Config.Message+passwd, false)); err != nil {
		return fmt.Errorf("send Slack DM to a created user(Slack User ID: %s): %w", user.ID, err)
	}
	return nil
}
