package handlers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"go-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var inactivityTimer *time.Timer

func ResetInactivityTimer(bot *tgbotapi.BotAPI, chatID int64) {
	if inactivityTimer != nil {
		inactivityTimer.Stop()
	}
	inactivityTimer = time.AfterFunc(3*time.Minute, func() {
		SendPasswordRequest(bot, chatID)
	})
}

func SendPasswordRequest(bot *tgbotapi.BotAPI, chatID int64) {
	currentAction = "password_request"
	msg := tgbotapi.NewMessage(chatID, "<b>Введите пароль для доступа к настройкам</b>")
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func SendMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Комплексное обновление", "update_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Остановить все задачи", "stop_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Steam", "steam")),
	)

	if messageID == 0 {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = keyboard
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	} else {
		editMessage := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, keyboard)
		editMessage.ParseMode = "HTML"
		if _, err := bot.Send(editMessage); err != nil {
			log.Printf("Failed to edit message: %v", err)
		}
	}
}

func SendAccountManagementMenu(bot *tgbotapi.BotAPI, chatID int64) {
	var buttons [][]tgbotapi.InlineKeyboardButton

	for id, account := range config.GetConfig().Accounts {
		buttonText := account.Username
		if strings.HasPrefix(account.Username, "disabled_") {
			buttonText = "🔒 " + account.Username[len("disabled_"):]
		}
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, "manage_account_"+strconv.Itoa(id)),
		))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_password"),
	))

	response := "<b>Управление аккаунтами</b>\n\nВыберите аккаунт для управления."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	ResetInactivityTimer(bot, chatID)
}
