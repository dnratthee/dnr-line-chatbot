// Line Messaging API Webhook

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func flexMessage(jsonBytes []byte) messaging_api.FlexComponentInterface {
	flexContainer, err := messaging_api.UnmarshalFlexContainer(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}
	return flexContainer
}

func init() {
	if os.Getenv("LINE_CHANNEL_SECRET") == "" {
		log.Fatal("LINE_CHANNEL_SECRET must be set")
	}
	if os.Getenv("LINE_CHANNEL_TOKEN") == "" {
		log.Fatal("LINE_CHANNEL_TOKEN must be set")
	}
}

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

					if strings.ToLower(message.Text) == "contact" {
						_, err = bot.ReplyMessage(
							&messaging_api.ReplyMessageRequest{
								ReplyToken: e.ReplyToken,
								Messages: []messaging_api.MessageInterface{
									&messaging_api.FlexMessage{
										AltText: "My Contact Information",
										Contents: flexMessage([]byte(`{
											"type": "bubble",
											"hero": {
												"type": "image",
												"url": "https://dnratthee.me/images/DNR.png",
												"size": "full",
												"aspectRatio": "4:5",
												"action": {
												"type": "uri",
												"uri": "https://dnratthee.me",
												"label": "DNRatthee"
												},
												"aspectMode": "cover",
												"backgroundColor": "#00cc99",
												"position": "relative",
												"margin": "none"
											},
											"body": {
												"type": "box",
												"layout": "vertical",
												"contents": [
													{
														"type": "text",
														"text": "Ratthee Jarathbenjawong",
														"weight": "bold",
														"size": "lg"
													}
												]
											},
											"footer": {
												"type": "box",
												"layout": "vertical",
												"spacing": "sm",
												"contents": [
													{
														"type": "button",
														"style": "link",
														"height": "sm",
														"action": {
														"type": "uri",
														"label": "E-MAIL",
														"uri": "mailto:info@dnratthee.me"
														}
													},
													{
														"type": "button",
														"style": "link",
														"height": "sm",
														"action": {
														"type": "uri",
														"label": "WEBSITE",
														"uri": "https://dnratthee.me"
														}
													}
												],
												"flex": 0
											}
										}`)),
									},
								},
							},
						)
						if err != nil {
							log.Print(err)
						}
						return
					}

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
