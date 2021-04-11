package controller

import (
	"context"
	"errors"

	"github.com/slack-go/slack"
)

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
		var err error
		userPagination, err = userPagination.Next(ctx)
		if userPagination.Done(err) {
			return slack.User{}, errors.New("user isn't found: " + name)
		}
	}
}
