package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sey_name = [4]string{"вито скалетто", "саша", "виктория", "вика"}

var bot *tgbotapi.BotAPI

var ChatID int64
var answers = []string{
	"Здарова",
	"Это пригожин женя",
	"Ну ты конечно даешь",
	"Как твои дела",
	"Хаха",
	"Звучит как анекдот",
}

func connect() {
	var err error
	if bot, err = tgbotapi.NewBotAPI("5109186321:AAEiyfx1qz9QYK9gxJ6lDrhATtAk5W-ZrIo"); err != nil {
		panic("Connection error")
	}

}

func MakeRequest() string {
	var s string
	resp, err := http.Get("https://api.adviceslip.com/advice")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	s = string(body)
	s = s[32 : len(s)-3]
	if s[0] == '"' {
		return s[1:]
	} else {
		return s
	}

}

func sendMessage(msg string) {
	msgConfig := tgbotapi.NewMessage(ChatID, msg)
	bot.Send(msgConfig)
	log.Println("bots: ", msgConfig)
}

func getFort() string {
	ind := rand.Intn(len(answers))
	return answers[ind]
}

func sendAnswer(update *tgbotapi.Update) {
	var get string
	get = getFort()
	msg := tgbotapi.NewMessage(ChatID, get)
	msg.ReplyToMessageID = int(update.Message.MessageID)
	bot.Send(msg)
	log.Println("bots: ", get)
}

func sendAnswerWrong(update *tgbotapi.Update) {
	var get string
	get = MakeRequest()
	msg := tgbotapi.NewMessage(ChatID, get)
	msg.ReplyToMessageID = int(update.Message.MessageID)
	bot.Send(msg)
	log.Println("bots: ", get)
}

func updateText(update *tgbotapi.Update) bool {

	if update.Message != nil && update.Message.Text == "/" {
		return false
	}

	msgLower := strings.ToLower(update.Message.Text)
	for _, name := range sey_name {
		if strings.Contains(msgLower, name) {
			return true
		}
	}
	return false
}

func main() {
	connect()
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 30
	updates := bot.GetUpdatesChan(upd)

	for updf := range updates {
		ChatID = updf.Message.Chat.ID
		log.Println("user: ", updf.Message.Text)

		if updf.Message != nil {
			rand.Seed(time.Now().Unix())
			var iii = rand.Intn(4)
			if iii == 1 {
				sendMessage("Да")
			} else if iii == 2 {
				sendMessage("Нет")
			} else if iii == 3 {
				sendMessage("Не знаю")
			} else {
				sendAnswer(&updf)
			}

		}

		if updf.Message != nil && updf.Message.Text == "/start" {
			sendMessage("Задавай свои вопросы")
		}
		if updf.Message.Text == "Ты писька" || updf.Message.Text == "Ты жопа" {
			sendMessage("Нет ты")
		}

		if updateText(&updf) {
			sendAnswer(&updf)
		}

		if updf.Message != nil && updf.Message.Text == "/inter" {
			sendAnswerWrong(&updf)
		}

	}

}
