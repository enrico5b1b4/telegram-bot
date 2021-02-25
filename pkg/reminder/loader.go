package reminder

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

import (
	"fmt"

	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/pkg/cron"
	"github.com/enrico5b1b4/telegram-bot/pkg/telegram"
)

type LoaderServicer interface {
	LoadSchedulesFromDB() (int, error)
	ReloadSchedulesForChat(chatID int) (int, error)
}

type LoaderService struct {
	b                   telegram.TBWrapBot
	scheduler           cron.Scheduler
	reminderStore       Storer
	reminderJobService  CronFuncServicer
	chatPreferenceStore chatpreference.Storer
}

func NewLoaderService(
	b telegram.TBWrapBot,
	scheduler cron.Scheduler,
	reminderStore Storer,
	chatPreferenceStore chatpreference.Storer,
	reminderJobService CronFuncServicer,
) *LoaderService {
	return &LoaderService{
		b:                   b,
		scheduler:           scheduler,
		reminderStore:       reminderStore,
		chatPreferenceStore: chatPreferenceStore,
		reminderJobService:  reminderJobService,
	}
}

// LoadSchedulesFromDB loads reminders from the DB
// and creates schedules on the scheduler.
// Only Active reminders will have a schedule created
func (s *LoaderService) LoadSchedulesFromDB() (int, error) {
	remindersAdded := 0
	rmdrListByChat, err := s.reminderStore.GetAllRemindersByChat()
	if err != nil {
		return 0, err
	}

	for chatID := range rmdrListByChat {
		chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
		if err != nil {
			return 0, err
		}

		for i := range rmdrListByChat[chatID] {
			if rmdrListByChat[chatID][i].Status != cron.Active {
				continue
			}

			schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rmdrListByChat[chatID][i].Job.Schedule)
			reminderCronID, err := s.scheduler.Add(
				schedule,
				NewCronFunc(s.reminderJobService, s.b, &rmdrListByChat[chatID][i]),
			)
			if err != nil {
				return 0, err
			}

			rmdrListByChat[chatID][i].CronID = reminderCronID
			err = s.reminderStore.UpdateReminder(&rmdrListByChat[chatID][i])
			if err != nil {
				return 0, err
			}

			remindersAdded++
		}
	}

	return len(rmdrListByChat), nil
}

// ReloadSchedulesForChat reschedules reminders for a particular chat.
// If a schedule is already present on the scheduler it is removed before being added again
// This is needed as the timezone of the chat might have changed
func (s *LoaderService) ReloadSchedulesForChat(chatID int) (int, error) {
	remindersLoaded := 0
	rmdrListByChat, err := s.reminderStore.GetAllRemindersByChatID(chatID)
	if err != nil {
		return 0, err
	}

	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return 0, err
	}

	for i := range rmdrListByChat {
		if entry := s.scheduler.GetEntryByID(rmdrListByChat[i].CronID); entry.ID != 0 {
			s.scheduler.Remove(entry.ID)
		}

		if rmdrListByChat[i].Status != cron.Active {
			continue
		}

		schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rmdrListByChat[i].Job.Schedule)
		reminderID, err := s.scheduler.Add(
			schedule,
			NewCronFunc(s.reminderJobService, s.b, &rmdrListByChat[i]),
		)
		if err != nil {
			return 0, err
		}

		rmdrListByChat[i].CronID = reminderID
		err = s.reminderStore.UpdateReminder(&rmdrListByChat[i])
		if err != nil {
			return 0, err
		}

		remindersLoaded++
	}

	return len(rmdrListByChat), nil
}
