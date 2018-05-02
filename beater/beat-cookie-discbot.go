package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/imiMoisesEducation/beat-cookie-discbot/config"
	bot "github.com/imiMoisesEducation/beat-cookie-discbot/discord"
)

type Beatcookiediscbot struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {

	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	err := bot.SetToken("NDQwOTE1NjMzMDIyOTU5NjE2.Dcp6Gg.06PrjRhgC3eFhqr4l8H81x1edqw")

	if err != nil {
		panic(fmt.Sprintf("Couldnt connect to discord, reason: %v\n", err.Error()))
	}

	bt := &Beatcookiediscbot{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Beatcookiediscbot) Run(b *beat.Beat) error {
	logp.Info("beat-cookie-discbot is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	fmt.Printf("\nboop, beat\n")
	if err != nil {
		fmt.Println(err)
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	fmt.Printf("boop, beat4")
	for {
		select {
		case <-bt.done:
			fmt.Printf("boop, beat5")
			return nil
		case <-ticker.C:
			fmt.Printf("\nboop, beat6\n")
		}

		ch := bot.GetLatestMessages()

		if ch != nil {
			fmt.Printf("boop, beat 2")
			for message := range ch {
				fmt.Printf("boop, beat 3")
				tm, err := message.Timestamp.Parse()

				if err != nil {
					continue
				}

				reactions := []common.MapStr{}

				for _, reaction := range message.Reactions {
					reactions = append(reactions, common.MapStr{
						"count":     reaction.Count,
						"emojiID":   reaction.Emoji.ID,
						"emojiName": reaction.Emoji.Name,
					})
				}

				mentions := []string{}

				for _, mention := range message.Mentions {
					mentions = append(mentions, mention.ID)
				}

				event := common.MapStr{
					"attatchmentCount":   len(message.Attachments),
					"authorID":           message.Author.ID,
					"channelID":          message.ChannelID,
					"content":            message.ContentWithMentionsReplaced(),
					"embedsCount":        len(message.Embeds),
					"messageID":          message.ID,
					"didMentionEveryone": message.MentionEveryone,
					"mentionRoles":       message.MentionRoles,
					"mentionsCount":      message.Mentions,
					"mentions":           mentions,
					"reactions":          reactions,
					"type":               message.Type,
				}
				beet := beat.Event{Timestamp: tm, Fields: event}
				bt.client.Publish(beet)
			}
		}
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Beatcookiediscbot) Stop() {
	bt.client.Close()
	close(bt.done)
	bot.TurnOff()
}
