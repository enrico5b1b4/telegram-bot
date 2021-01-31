package command

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
)

// nolint:lll
const HandlePatternGetTimezone = `/gettimezone`

func HandleGetTimezone(store chatpreference.Storer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		cp, err := store.GetChatPreference(int(c.ChatID()))
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Your timezone is: %s", cp.TimeZone))

		return err
	}
}
