package command

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

import (
	"errors"

	"github.com/enrico5b1b4/telegram-bot/pkg/cron"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
)

type RemindDeleteServicer interface {
	DeleteReminder(chatID, ID int) error
}

type RemindDeleteService struct {
	reminderStore reminder.Storer
	scheduler     cron.Scheduler
}

func NewRemindeDeleteService(reminderStore reminder.Storer, scheduler cron.Scheduler) *RemindDeleteService {
	return &RemindDeleteService{
		reminderStore: reminderStore,
		scheduler:     scheduler,
	}
}

func (s *RemindDeleteService) DeleteReminder(chatID, id int) error {
	r, err := s.reminderStore.GetReminder(chatID, id)
	if err != nil {
		return err
	}

	if chatID != r.ChatID {
		return errors.New("unauthorised to delete reminder")
	}

	s.scheduler.Remove(r.CronID)

	return s.reminderStore.DeleteReminder(r.ChatID, id)
}
