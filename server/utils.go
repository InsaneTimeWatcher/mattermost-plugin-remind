package main

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nicksnyder/go-i18n/i18n"
)

func MaxTime(a time.Time, b time.Time) time.Time {
	if (a.Before(b)) {
		return b
	} else {
		return a
	}
}

func (p *Plugin) translation(user *model.User) (i18n.TranslateFunc, string) {
	locale := "en"
	for l := range p.locales {
		if user.Locale == l {
			locale = user.Locale
			break
		}
	}
	return p.GetUserTranslations(locale), locale
}

func (p *Plugin) location(user *model.User) *time.Location {
	tz := user.GetPreferredTimezone()
	if tz == "" {
		// Use server's timezone
		return time.Now().Location()
	} else {
		location, _ := time.LoadLocation(tz)
		return location
	}

}
