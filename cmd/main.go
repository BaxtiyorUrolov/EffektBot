package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"effektbot/admin"
	"effektbot/state"
	"effektbot/storage"
	"effektbot/effekt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	connStr := "user=godb password=0208 dbname=effektbot sslmode=disable"
	db, err := storage.OpenDatabase(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	botToken := "6902655696:AAEtKAL78CG86DhjAYb-QVQrTVAGysTpLDA"
	botInstance, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	offset := 0
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot...")
			return
		default:
			updates, err := botInstance.GetUpdates(tgbotapi.NewUpdate(offset))
			if err != nil {
				log.Printf("Error getting updates: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			for _, update := range updates {
				handleUpdate(update, db, botInstance)
				offset = update.UpdateID + 1
			}
		}
	}
}

func handleUpdate(update tgbotapi.Update, db *sql.DB, botInstance *tgbotapi.BotAPI) {
	if update.Message != nil {
		handleMessage(update.Message, db, botInstance)
	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update.CallbackQuery, db, botInstance)
	} else {
		log.Printf("Unsupported update type: %T", update)
	}
}

func handleMessage(msg *tgbotapi.Message, db *sql.DB, botInstance *tgbotapi.BotAPI) {
	chatID := msg.Chat.ID
	text := msg.Text

	log.Printf("Received message: %s", text)

	if userState, exists := state.UserStates[chatID]; exists {
		if strings.HasPrefix(userState, "waiting_for_photo_") {
			effectID := strings.TrimPrefix(userState, "waiting_for_photo_")
			effekt.HandleUserPhoto(msg, effectID, botInstance)
			delete(state.UserStates, chatID)
			return
		}

		switch userState {
		case "waiting_for_broadcast_message":
			admin.HandleBroadcastMessage(msg, db, botInstance)
			delete(state.UserStates, chatID)
			return
		case "waiting_for_channel_link":
			admin.HandleChannelLink(msg, db, botInstance)
			delete(state.UserStates, chatID)
			return
		case "waiting_for_admin_id":
			admin.HandleAdminAdd(msg, db, botInstance)
			delete(state.UserStates, chatID)
			return
		case "waiting_for_admin_id_remove":
			admin.HandleAdminRemove(msg, db, botInstance)
			delete(state.UserStates, chatID)
			return
		}
	}

	if text == "/start" {
		handleStartCommand(msg, db, botInstance)
		storage.AddUserToDatabase(db, int(msg.Chat.ID))
	} else if text == "/admin" {
		admin.HandleAdminCommand(msg, db, botInstance)
	} else {
		handleDefaultMessage(msg, db, botInstance)
	}
}

func handleStartCommand(msg *tgbotapi.Message, db *sql.DB, botInstance *tgbotapi.BotAPI) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	log.Printf("Adding user to database: %d ", userID)
	err := storage.AddUserToDatabase(db, userID)
	if err != nil {
		log.Printf("Error adding user to database: %v", err)
		return
	}

	channels, err := storage.GetChannelsFromDatabase(db)
	if err != nil {
		log.Printf("Error getting channels from database: %v", err)
		return
	}

	log.Printf("Checking subscription for user %d", chatID)
	if isUserSubscribedToChannels(chatID, channels, botInstance) {
		msg := tgbotapi.NewMessage(chatID, "🇺🇿 UZ Yangiliklar @MRC_groupuz \n\n Assalomu Alaykum kerakli boʻlimni tanlang ✅ \n\n➖➖➖➖➖➖➖➖➖➖➖➖➖➖➖➖ \n\n 🇷🇺 RUS Новости @MRC_groupuz \n\n Здравствуйте, Выберите раздел в соответствии  ✅ \n\n admin: @MRC_Admin")
		effectButton := tgbotapi.NewInlineKeyboardButtonData("🎆 Rasmga effekt", "apply_effect_19914101")
		newsButton := tgbotapi.NewInlineKeyboardButtonData("🆕 Yangiliklar", "news")
		var inlineKeyboard tgbotapi.InlineKeyboardMarkup

		if storage.IsAdmin(int(chatID), db) {
			adminButton := tgbotapi.NewInlineKeyboardButtonData("Admin", "admin")
			inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(effectButton, newsButton),
				tgbotapi.NewInlineKeyboardRow(adminButton),
			)
		} else {
			inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(effectButton, newsButton),
			)
		}

		msg.ReplyMarkup = inlineKeyboard
		botInstance.Send(msg)
	} else {
		log.Printf("User %d is not subscribed to required channels", chatID)
		inlineKeyboard := createSubscriptionKeyboard(channels)
		msg := tgbotapi.NewMessage(chatID, "Iltimos, avval kanallarga azo bo'ling:")
		msg.ReplyMarkup = inlineKeyboard
		botInstance.Send(msg)
	}
}

func handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, db *sql.DB, botInstance *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	channels, err := storage.GetChannelsFromDatabase(db)
	if err != nil {
		log.Printf("Error getting channels from database: %v", err)
		return
	}

	if callbackQuery.Data == "check_subscription" {
		if isUserSubscribedToChannels(chatID, channels, botInstance) {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)

			msg := tgbotapi.NewMessage(chatID, "Assalomu alaykum, siz kanallarga azo bo'ldingiz!")
				effectButton := tgbotapi.NewInlineKeyboardButtonData("🎆 Rasmga effekt", "apply_effect_19914101")
				newsButton := tgbotapi.NewInlineKeyboardButtonData("🆕 Yangiliklar", "news")
				var inlineKeyboard tgbotapi.InlineKeyboardMarkup

				if storage.IsAdmin(int(chatID), db) {
					adminButton := tgbotapi.NewInlineKeyboardButtonData("Admin", "admin")
					inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(effectButton, newsButton),
						tgbotapi.NewInlineKeyboardRow(adminButton),
					)
				} else {
					inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(effectButton, newsButton),
					)
				}
			msg.ReplyMarkup = inlineKeyboard
			botInstance.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Iltimos, kanallarga azo bo'ling.")
			inlineKeyboard := createSubscriptionKeyboard(channels)
			msg.ReplyMarkup = inlineKeyboard
			botInstance.Send(msg)
		}
	} else if strings.HasPrefix(callbackQuery.Data, "apply_effect_") {
		effectID := strings.TrimPrefix(callbackQuery.Data, "apply_effect_")
		effekt.HandleApplyEffect(callbackQuery, effectID, botInstance)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if strings.HasPrefix(callbackQuery.Data, "send_photo_") {
		effectID := strings.TrimPrefix(callbackQuery.Data, "send_photo_")
		handleSendPhotoRequest(callbackQuery, effectID, botInstance)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if strings.HasPrefix(callbackQuery.Data, "next_effect_") {
		nextEffectID := strings.TrimPrefix(callbackQuery.Data, "next_effect_")
		effekt.HandleNextEffect(callbackQuery, nextEffectID, botInstance)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if strings.HasPrefix(callbackQuery.Data, "previous_effect_") {
		previousEffectID := strings.TrimPrefix(callbackQuery.Data, "previous_effect_")
		effekt.HandlePreviousEffect(callbackQuery, previousEffectID, botInstance)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if strings.HasPrefix(callbackQuery.Data, "delete_channel_") {
		channel := strings.TrimPrefix(callbackQuery.Data, "delete_channel_")
		admin.AskForChannelDeletionConfirmation(chatID, messageID, channel, botInstance)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if strings.HasPrefix(callbackQuery.Data, "confirm_delete_channel_") {
		channel := strings.TrimPrefix(callbackQuery.Data, "confirm_delete_channel_")
		admin.DeleteChannel(chatID, messageID, channel, db, botInstance)
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
			botInstance.Send(deleteMsg)
	} else if callbackQuery.Data == "cancel_delete_channel" {
		admin.CancelChannelDeletion(chatID, messageID, botInstance)
	}
}

func handleSendPhotoRequest(callbackQuery *tgbotapi.CallbackQuery, effectID string, botInstance *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Iltimos, rasm yuboring:")
	state.UserStates[chatID] = "waiting_for_photo_" + effectID
	botInstance.Send(msg)
}

func handleDefaultMessage(msg *tgbotapi.Message, db *sql.DB, botInstance *tgbotapi.BotAPI) {
	chatID := msg.Chat.ID
	text := msg.Text

	switch text {
	case "Kanal qo'shish":
		state.UserStates[chatID] = "waiting_for_channel_link"
		msgResponse := tgbotapi.NewMessage(chatID, "Kanal linkini yuboring (masalan, https://t.me/your_channel):")
		botInstance.Send(msgResponse)
	case "Admin qo'shish":
		state.UserStates[chatID] = "waiting_for_admin_id"
		msgResponse := tgbotapi.NewMessage(chatID, "Iltimos, yangi admin ID sini yuboring:")
		botInstance.Send(msgResponse)
	case "Admin o'chirish":
		state.UserStates[chatID] = "waiting_for_admin_id_remove"
		msgResponse := tgbotapi.NewMessage(chatID, "Iltimos, admin ID sini o'chirish uchun yuboring:")
		botInstance.Send(msgResponse)
	case "Kanal o'chirish":
		admin.DisplayChannelsForDeletion(chatID, db, botInstance)
	case "Statistika":
		admin.HandleStatistics(msg, db, botInstance)
	case "Habar yuborish":
		state.UserStates[chatID] = "waiting_for_broadcast_message"
		msgResponse := tgbotapi.NewMessage(chatID, "Iltimos, yubormoqchi bo'lgan habaringizni kiriting (Bekor qilish uchun /cancel):")
		botInstance.Send(msgResponse)
	default:
		msgResponse := tgbotapi.NewMessage(chatID, "Har qanday boshqa xabarlarni shu yerda ko'rib chiqish mumkin")
		botInstance.Send(msgResponse)
	}
}

func isUserSubscribedToChannels(chatID int64, channels []string, botInstance *tgbotapi.BotAPI) bool {
	for _, channel := range channels {
		log.Printf("Checking subscription to channel: %s", channel)
		chat, err := botInstance.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: "@" + channel})
		if err != nil {
			log.Printf("Error getting chat info for channel %s: %v", channel, err)
			return false
		}

		member, err := botInstance.GetChatMember(tgbotapi.ChatConfigWithUser{
			ChatID: chat.ID,
			UserID: int(chatID),
		})
		if err != nil {
			log.Printf("Error getting chat member info for channel %s: %v", channel, err)
			return false
		}
		if member.Status == "left" || member.Status == "kicked" {
			log.Printf("User %d is not subscribed to channel %s", chatID, channel)
			return false
		}
	}
	return true
}

func createSubscriptionKeyboard(channels []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, channel := range channels {
		channelName := strings.TrimPrefix(channel, "https://t.me/")
		button := tgbotapi.NewInlineKeyboardButtonURL(channelName, "https://t.me/"+channelName)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	checkButton := tgbotapi.NewInlineKeyboardButtonData("Azo bo'ldim", "check_subscription")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(checkButton))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
