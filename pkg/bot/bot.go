package bot

import (
	"log"

	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/gettimezone"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindat"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/reminddaymonth"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/reminddayofweek"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/reminddelete"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/reminddetail"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindevery"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindeveryday"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindeverydaynumber"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindeverydaynumbermonth"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindeverydayofweek"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindhelp"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindin"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindlist"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/remindwhen"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/settimezone"
	"github.com/enrico5b1b4/telegram-bot/pkg/cron"
	"github.com/enrico5b1b4/telegram-bot/pkg/date"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
	"github.com/enrico5b1b4/telegram-bot/pkg/telegram"
	"go.etcd.io/bbolt"
)

type Bot struct {
	cronScheduler cron.Scheduler
	telegramBot   telegram.TBWrapBot
}

// nolint:funlen
func New(
	allowedChats []int,
	database *bbolt.DB,
	telegramBot telegram.TBWrapBot,
) *Bot {
	cronScheduler := cron.NewScheduler()
	reminderStore := reminder.NewStore(database)
	chatPreferenceStore := chatpreference.NewStore(database)
	chatPreferenceService := chatpreference.NewService(chatPreferenceStore)
	remindCronFuncService := reminder.NewCronFuncService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore)
	remindListService := remindlist.NewService(reminderStore, cronScheduler, chatPreferenceStore)
	remindDeleteService := reminddelete.NewService(reminderStore, cronScheduler)
	reminderScheduler := reminder.NewScheduler(telegramBot, remindCronFuncService, reminderStore, cronScheduler, chatPreferenceStore)
	remindDateService := reminder.NewService(reminderScheduler, reminderStore, chatPreferenceStore, date.RealTimeNow)
	remindDetailService := reminddetail.NewService(reminderStore, cronScheduler, chatPreferenceStore)
	reminderLoader := reminder.NewLoaderService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore, remindCronFuncService)
	setTimeZoneService := settimezone.NewService(chatPreferenceStore, reminderLoader)
	remindDetailButtons := reminddetail.NewButtons()
	remindListButtons := remindlist.NewButtons()
	reminderCompleteButtons := reminder.NewButtons()

	chatPreferenceService.CreateDefaultChatPreferences(allowedChats)

	// check if DB exists and load schedules
	remindersLoaded, err := reminderLoader.LoadSchedulesFromDB()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d reminders", remindersLoaded)

	telegramBot.Handle(remindlist.HandlePattern, remindlist.HandleRemindList(remindListService, remindListButtons))
	telegramBot.Handle(remindhelp.HandlePattern, remindhelp.HandleRemindHelp())
	telegramBot.HandleMultiRegExp(reminddetail.HandlePattern, reminddetail.HandleRemindDetail(remindDetailService, reminddetail.NewButtons()))
	telegramBot.HandleMultiRegExp(reminddelete.HandlePattern, reminddelete.HandleRemindDelete(remindDeleteService))
	telegramBot.HandleRegExp(
		reminddaymonth.HandlePattern,
		reminddaymonth.HandleRemindDayMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		reminddayofweek.HandlePattern,
		reminddayofweek.HandleRemindDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumber.HandlePattern,
		remindeverydaynumber.HandleRemindEveryDayNumber(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumbermonth.HandlePattern,
		remindeverydaynumbermonth.HandleRemindEveryDayNumberMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindin.HandlePattern,
		remindin.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindevery.HandlePattern,
		remindevery.HandleRemindEvery(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindwhen.HandlePattern,
		remindwhen.HandleRemindWhen(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydayofweek.HandlePattern,
		remindeverydayofweek.HandleRemindEveryDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeveryday.HandlePattern,
		remindeveryday.HandleRemindEveryDay(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindat.HandlePattern,
		remindat.HandleRemindAt(remindDateService),
	)
	telegramBot.Handle(gettimezone.HandlePattern, gettimezone.HandleGetTimezone(chatPreferenceStore))
	telegramBot.HandleRegExp(settimezone.HandlePattern, settimezone.HandleSetTimezone(setTimeZoneService))

	// buttons
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailCloseCommandBtn],
		reminddetail.HandleCloseBtn(),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailDeleteBtn],
		reminddetail.HandleReminderDetailDeleteBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailShowReminderCommandBtn],
		reminddetail.HandleReminderShowReminderCommandBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListRemoveCompletedRemindersBtn],
		remindlist.HandleReminderListRemoveCompletedRemindersBtn(remindListService),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListCloseCommandBtn],
		remindlist.HandleCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze10MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 10}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze20MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 20}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze30MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 30}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze1HourBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 60}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeThisAfternoonBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeThisEveningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowMorningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   9,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowAfternoonBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowEveningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeBtn],
		reminder.HandleReminderSnoozeBtn(reminderStore),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeCloseBtn],
		reminder.HandleReminderSnoozeCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.CompleteBtn],
		reminder.HandleReminderCompleteBtn(remindCronFuncService, reminderStore),
	)

	return &Bot{
		cronScheduler: cronScheduler,
		telegramBot:   telegramBot,
	}
}

func (b *Bot) Start() {
	b.cronScheduler.Start()
	b.telegramBot.Start()
}
