package handlers

import (
	"log"
	"strconv"
	"strings"

	"go-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tempAccountID int
var currentAction string

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if currentAction != "" {
		switch currentAction {
		case "edit_username":
			account := config.GetConfig().Accounts[tempAccountID]
			account.Username = message.Text
			config.GetConfig().Accounts[tempAccountID] = account
			config.SaveConfig("config.json")
			SendAccountManagementMenu(bot, message.Chat.ID)
		case "edit_password":
			account := config.GetConfig().Accounts[tempAccountID]
			account.Password = message.Text
			config.GetConfig().Accounts[tempAccountID] = account
			config.SaveConfig("config.json")
			SendAccountManagementMenu(bot, message.Chat.ID)
		case "password_request":
			if message.Text == config.GetConfig().AdminPassword {
				SendAccountManagementMenu(bot, message.Chat.ID)
			} else {
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Неверный пароль. Попробуйте снова."))
				SendPasswordRequest(bot, message.Chat.ID)
			}
		}
		currentAction = ""
		return
	}

	if message.Text == "/start" {
		SendMainMenu(bot, message.Chat.ID, 0, "<b>Главное меню</b>\n\nВыберите действие:")
	} else {
		// Ваш код для обработки других сообщений
		log.Printf("Received a message: %s", message.Text)
	}
}

func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	callbackQueryData := callbackQuery.Data
	var response string
	var keyboard tgbotapi.InlineKeyboardMarkup

	switch callbackQueryData {
	case "update_all":
		response = "<b>Будет выполнено комплексное обновление. Подтвердите действие.</b>"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да", "confirm_update_all")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет", "cancel")),
		)
	case "stop_all":
		response = "<b>Все задачи будут остановлены. Подтвердите действие.</b>"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да", "confirm_stop_all")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет", "cancel")),
		)
	case "steam":
		response = "<b>Меню: Steam</b>\n\nВыберите действие:"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Обновление игр Steam", "update_steam_games")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Установка игр Steam", "install_steam_games")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Настройка аккаунтов Steam", "manage_steam_accounts")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", "back")),
		)
	case "update_steam_games":
		response = "<b>Меню: Обновление игр Steam</b>\n\nПодтвердите действие."
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да", "confirm_update_steam_games")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет", "cancel")),
		)
	case "install_steam_games":
		response = "<b>Меню: Установка игр Steam</b>\n\nПодтвердите действие."
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да", "confirm_install_steam_games")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет", "cancel")),
		)
	case "manage_steam_accounts":
		SendPasswordRequest(bot, chatID)
		return
	case "confirm_update_all":
		response = "<b>Комплексное обновление выполняется...</b>"
		SendMainMenu(bot, chatID, messageID, response)
		return
	case "confirm_stop_all":
		response = "<b>Остановка всех задач выполняется...</b>"
		SendMainMenu(bot, chatID, messageID, response)
		return
	case "confirm_update_steam_games":
		response = "<b>Игры Steam обновляются...</b>"
		SendMainMenu(bot, chatID, messageID, response)
		return
	case "confirm_install_steam_games":
		response = "<b>Игры Steam устанавливаются...</b>"
		SendMainMenu(bot, chatID, messageID, response)
		return
	case "cancel":
		SendMainMenu(bot, chatID, messageID, "<b>Главное меню</b>\n\nВыберите действие:")
		return
	case "back":
		SendMainMenu(bot, chatID, messageID, "<b>Главное меню</b>\n\nВыберите действие:")
		return
	case "back_to_password":
		SendPasswordRequest(bot, chatID)
		return

	default:
		// Разбор команд для управления аккаунтами
		if strings.HasPrefix(callbackQueryData, "manage_account_") {
			accountID, err := strconv.Atoi(callbackQueryData[15:])
			if err == nil {
				ManageAccount(bot, chatID, accountID)
			}
		} else if strings.HasPrefix(callbackQueryData, "edit_username_") {
			accountID, err := strconv.Atoi(callbackQueryData[14:])
			if err == nil {
				tempAccountID = accountID
				currentAction = "edit_username"
				msg := tgbotapi.NewMessage(chatID, "Введите новый логин:")
				bot.Send(msg)
			}
		} else if strings.HasPrefix(callbackQueryData, "edit_password_") {
			accountID, err := strconv.Atoi(callbackQueryData[14:])
			if err == nil {
				tempAccountID = accountID
				currentAction = "edit_password"
				msg := tgbotapi.NewMessage(chatID, "Введите новый пароль:")
				bot.Send(msg)
			}
		} else if strings.HasPrefix(callbackQueryData, "delete_account_") {
			accountID, err := strconv.Atoi(callbackQueryData[15:])
			if err == nil {
				delete(config.GetConfig().Accounts, accountID)
				config.SaveConfig("config.json")
				SendAccountManagementMenu(bot, chatID)
			}
		} else if strings.HasPrefix(callbackQueryData, "disable_account_") {
			accountID, err := strconv.Atoi(callbackQueryData[16:])
			if err == nil {
				account := config.GetConfig().Accounts[accountID]
				if strings.HasPrefix(account.Username, "disabled_") {
					account.Username = account.Username[len("disabled_"):]
				} else {
					account.Username = "disabled_" + account.Username
				}
				config.GetConfig().Accounts[accountID] = account
				config.SaveConfig("config.json")
				SendAccountManagementMenu(bot, chatID)
			}
		}
	}

	// Обновление сообщения для отображения ответа и скрытия кнопок
	editMessage := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, response, keyboard)
	editMessage.ParseMode = "HTML"
	if _, err := bot.Send(editMessage); err != nil {
		log.Printf("Failed to send edited message: %v", err)
	}

	// Ответ на CallbackQuery, чтобы скрыть индикатор загрузки
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Printf("Failed to send callback response: %v", err)
	}
}

func ManageAccount(bot *tgbotapi.BotAPI, chatID int64, accountID int) {
	account := config.GetConfig().Accounts[accountID]
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить логин", "edit_username_"+strconv.Itoa(accountID)),
			tgbotapi.NewInlineKeyboardButtonData("Изменить пароль", "edit_password_"+strconv.Itoa(accountID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить аккаунт", "delete_account_"+strconv.Itoa(accountID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_account_management"),
		),
	)

	response := "<b>Управление аккаунтом:</b>\n\nВыберите действие для аккаунта <b>" + account.Username + "</b>."
	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send account management message: %v", err)
	}

	ResetInactivityTimer(bot, chatID)
}
