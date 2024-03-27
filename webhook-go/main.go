// Line Messaging API Webhook

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {

	handler, err := webhook.NewWebhookHandler(
		os.Getenv("LINE_CHANNEL_SECRET"),
	)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))

	handler.HandleEvents(func(req *webhook.CallbackRequest, r *http.Request) {
		if err != nil {
			log.Print(err)
			return
		}
		log.Println("Handling events...")
		for _, event := range req.Events {
			log.Printf("/callback called%+v...\n", event)
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					_, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								&messaging_api.TextMessage{
									Text: message.Text,
								},
							},
						},
					)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	http.Handle("/callback", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("http(s)://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
