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

// https://api.slack.com/methods/users.list
// > We recommend no more than 200 results at a time.
const usersOpitonLimit = 190

func (ctrl *Controller) GetSlackUser(ctx context.Context, name string) (slack.User, error) {
	userPagination := ctrl.SlackBot.GetUsersPaginated(slack.GetUsersOptionLimit(usersOpitonLimit))
	for {
		for _, user := range userPagination.Users {
			// https://api.slack.com/types/user
			if user.Profile.DisplayName == name {
				return user, nil
			}
		}
		userPagination, err := userPagination.Next(ctx)
		if userPagination.Done(err) {
			return slack.User{}, errors.New("user isn't found: " + name)
		}
	}
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
