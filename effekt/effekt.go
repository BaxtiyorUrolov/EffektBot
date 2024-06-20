package effekt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleApplyEffect(callbackQuery *tgbotapi.CallbackQuery, effectID string, botInstance *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID

	photoMsg := tgbotapi.NewPhotoShare(chatID, "https://t.me/photolabsuz/9")
	photoMsg.Caption = "Mana 4101. Keyingi tugmalarni ko'ring:"
	nextEffectID := "16897287"
	page := "1/7"
	effectButton := tgbotapi.NewInlineKeyboardButtonData("Tanlash", fmt.Sprintf("send_photo_%s", effectID))
	pageButton := tgbotapi.NewInlineKeyboardButtonData(page, fmt.Sprintf("%s/7", page))
	nextButton := tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("next_effect_%s", nextEffectID))
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(pageButton, nextButton),
		tgbotapi.NewInlineKeyboardRow(effectButton),
	)
	photoMsg.ReplyMarkup = inlineKeyboard

	botInstance.Send(photoMsg)
}


func HandleNextEffect(callbackQuery *tgbotapi.CallbackQuery, effectID string, botInstance *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID

	var photoURL string
	var caption string
	var nextEffectID string
	var previousEffectID string
	var page string
	var effectButtonData string

	switch effectID {
	case "16897287":
		photoURL = "https://t.me/photolabsuz/12"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "19990153"
		previousEffectID = "19914101"
		page = "2/7"
		effectButtonData = "send_photo_16897287"
	case "19990153":
		photoURL = "https://t.me/photolabsuz/4"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "39197678"
		previousEffectID = "16897287"
		page = "3/7"
		effectButtonData = "send_photo_19990153"
	case "39197678":
		photoURL = "https://t.me/photolabsuz/13"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "23059052"
		previousEffectID = "19990153"
		page = "4/7"
		effectButtonData = "send_photo_39197678"
	case "23059052":
		photoURL = "https://t.me/photolabsuz/14"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "2490"
		previousEffectID = "39197678"
		page = "5/7"
		effectButtonData = "send_photo_23059052"
	case "2490":
		photoURL = "https://t.me/photolabsuz/11"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "39175403"
		previousEffectID = "23059052"
		page = "6/7"
		effectButtonData = "send_photo_2490"
	case "39175403":
		photoURL = "https://t.me/photolabsuz/10"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "19914101"
		previousEffectID = "2490"
		page = "7/7"
		effectButtonData = "send_photo_39175403"
	}

	photoMsg := tgbotapi.NewPhotoShare(chatID, photoURL)
	photoMsg.Caption = caption
	effectButton := tgbotapi.NewInlineKeyboardButtonData("Tanlash", effectButtonData)
	nextButton := tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("next_effect_%s", nextEffectID))
	previousButton := tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("previous_effect_%s", previousEffectID))
	pageButton := tgbotapi.NewInlineKeyboardButtonData(page, fmt.Sprintf("%s/7", page))
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(previousButton, pageButton, nextButton),
		tgbotapi.NewInlineKeyboardRow(effectButton),
	)
	photoMsg.ReplyMarkup = inlineKeyboard

	botInstance.Send(photoMsg)
}

func HandlePreviousEffect(callbackQuery *tgbotapi.CallbackQuery, effectID string, botInstance *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID

	var photoURL string
	var caption string
	var nextEffectID string
	var previousEffectID string
	var page string
	var effectButtonData string

	switch effectID {
	case "19914101":
		photoURL = "https://t.me/photolabsuz/9"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "16897287"
		previousEffectID = "19914101"
		page = "1/7"
		effectButtonData = "send_photo_19914101"
	case "16897287":
		photoURL = "https://t.me/photolabsuz/12"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "19990153"
		previousEffectID = "19914101"
		page = "2/7"
		effectButtonData = "send_photo_16897287"
	case "19990153":
		photoURL = "https://t.me/photolabsuz/4"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "39197678"
		previousEffectID = "16897287"
		page = "3/7"
		effectButtonData = "send_photo_19990153"
	case "39197678":
		photoURL = "https://t.me/photolabsuz/13"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "23059052"
		previousEffectID = "19990153"
		page = "4/7"
		effectButtonData = "send_photo_39197678"
	case "23059052":
		photoURL = "https://t.me/photolabsuz/14"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "2490"
		previousEffectID = "39197678"
		page = "5/7"
		effectButtonData = "send_photo_23059052"
	case "2490":
		photoURL = "https://t.me/photolabsuz/11"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "39175403"
		previousEffectID = "23059052"
		page = "6/7"
		effectButtonData = "send_photo_2490"
	case "39175403":
		photoURL = "https://t.me/photolabsuz/10"
		caption = "Mana keyingi rasm. Keyingi tugmalarni ko'ring:"
		nextEffectID = "19914101"
		previousEffectID = "2490"
		page = "7/7"
		effectButtonData = "send_photo_39175403"
	}

	photoMsg := tgbotapi.NewPhotoShare(chatID, photoURL)
	photoMsg.Caption = caption
	effectButton := tgbotapi.NewInlineKeyboardButtonData("Tanlash", effectButtonData)
	nextButton := tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("next_effect_%s", nextEffectID))
	previousButton := tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("previous_effect_%s", previousEffectID))
	pageButton := tgbotapi.NewInlineKeyboardButtonData(page, fmt.Sprintf("%s/7", page))
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(previousButton, pageButton, nextButton),
		tgbotapi.NewInlineKeyboardRow(effectButton),
	)
	photoMsg.ReplyMarkup = inlineKeyboard

	botInstance.Send(photoMsg)
}

func HandleUserPhoto(msg *tgbotapi.Message, effectID string, botInstance *tgbotapi.BotAPI) {
	chatID := msg.Chat.ID

	if msg.Photo != nil {
		fileID := (*msg.Photo)[len(*msg.Photo)-1].FileID

		file, err := botInstance.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Printf("Error getting file: %v", err)
			botInstance.Send(tgbotapi.NewMessage(chatID, "Faylni olishda xatolik yuz berdi."))
			return
		}

		fileURL := file.Link(botInstance.Token)
		photoURL := fmt.Sprintf("http://195.2.84.169:8000/effect/?id=%s&photo=%s", effectID, fileURL)

		resp, err := http.Get(photoURL)
		if err != nil {
			log.Printf("Error sending request to API: %v", err)
			botInstance.Send(tgbotapi.NewMessage(chatID, "API ga so'rov yuborishda xatolik yuz berdi."))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			botInstance.Send(tgbotapi.NewMessage(chatID, "API dan javobni o'qishda xatolik yuz berdi."))
			return
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			botInstance.Send(tgbotapi.NewMessage(chatID, "API javobini tahlil qilishda xatolik yuz berdi."))
			return
		}

		if ok, found := result["result"].(map[string]interface{})["ok"].(bool); found && ok {
			imgURL := result["result"].(map[string]interface{})["img_url"].(string)
			photoMsg := tgbotapi.NewPhotoShare(chatID, imgURL)
			photoMsg.Caption = "Mana tayyorlangan rasm:"
			botInstance.Send(photoMsg)
		} else {
			botInstance.Send(tgbotapi.NewMessage(chatID, "Rasmni yaratishda xatolik yuz berdi."))
		}
	} else {
		botInstance.Send(tgbotapi.NewMessage(chatID, "Iltimos, rasm yuboring."))
	}
}
