package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const CommandTrigger = "remind"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		Username:         botName,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone or ~channel] [what] [when]",
		AutoCompleteDesc: "Set a reminder",
		DisplayName:      "Remind Plugin Command",
		Description:      "A command used to set a reminder",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {
	return p.API.UnregisterCommand(teamId, CommandTrigger)
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	user, uErr := p.API.GetUser(args.UserId)
	if uErr != nil {
		return &model.CommandResponse{}, uErr
	}
	T, locale := p.translation(user)
	location := p.location(user)

	if strings.HasSuffix(args.Command, T("help")) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(T("help.response")),
			Username:     botName,
		}, nil
	}

	if strings.HasSuffix(args.Command, T("list")) {
		listMessage := p.ListReminders(user, args.ChannelId)
		if listMessage != "" {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf(listMessage),
				Username:     botName,
			}, nil
		} else {
			return &model.CommandResponse{}, nil
		}

	}

	payload := strings.Trim(strings.Replace(args.Command, "/"+CommandTrigger, "", -1), " ")

	if strings.HasPrefix(payload, T("me")) ||
		strings.HasPrefix(payload, "@") ||
		strings.HasPrefix(payload, "~") {

		request := ReminderRequest{
			TeamId:   args.TeamId,
			Username: user.Username,
			Payload:  payload,
			Reminder: Reminder{},
		}
		response, err := p.ScheduleReminder(&request)

		if err != nil {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf(T("exception.response")),
				Username:     botName,
			}, nil
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("%s", response),
			Username:     botName,
		}, nil
	}

	// debug / troubleshooting commands

	// clear all reminders for current user
	if strings.HasSuffix(args.Command, "__clear") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(p.DeleteReminders(user)),
			Username:     botName,
		}, nil
	}
	// display the plugin version
	if strings.HasSuffix(args.Command, "__version") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(manifest.Version),
			Username:     botName,
		}, nil
	}

	// display the locale & location of user
	if strings.HasSuffix(args.Command, "__user") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(locale + " " + location.String()),
			Username:     botName,
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(T("exception.response")),
		Username:     botName,
	}, nil

}
