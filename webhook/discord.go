package webhook

import (
	"bytes"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"net/http"
)

type discordWebhookBody struct {
	Content string   `json:"content"`
	Embeds  []embeds `json:"embeds"`
}
type embeds struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

//{
//  "content": null,
//  "embeds": [
//    {
//      "title": "What's this?",
//      "description": "da",
//      "color": 65535
//    },
//    {
//      "title": "What's this?",
//      "description": "da",
//      "color": 65280
//    },
//    {
//      "color": 16711680
//    }
//  ]
//}

func DiscordError(v *bodyStruct.ErrorStruct) {

	body := discordWebhookBody{
		Embeds: []embeds{
			{
				Title:       "error",
				Description: "```json\n" + string(*utils.ToJsonIndentString(v)) + "```",
				Color:       16711680,
			},
		},
	}

	http.Post(config.Config.Discord.Webhook.Error, "application/json", bytes.NewReader(*utils.ToJsonString(body)))
}
func DiscordInfo(v *bodyStruct.ErrorStruct) {

	body := discordWebhookBody{
		Embeds: []embeds{
			{
				Title:       "info",
				Description: "```json\n" + string(*utils.ToJsonIndentString(v)) + "```",
				Color:       65535,
			},
		},
	}

	http.Post(config.Config.Discord.Webhook.Info, "application/json", bytes.NewReader(*utils.ToJsonString(body)))
}
