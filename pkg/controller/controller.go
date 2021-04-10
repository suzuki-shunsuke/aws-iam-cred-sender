package controller

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Controller struct {
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	Config    Config
	SlackBot  *slack.Client
	SlackUser *slack.Client
}

func New(ctx context.Context, param Param) (Controller, Param, error) {
	if param.LogLevel != "" {
		lvl, err := logrus.ParseLevel(param.LogLevel)
		if err != nil {
			return Controller{}, param, errors.New("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	return Controller{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, param, nil
}
