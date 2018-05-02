package discord

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	session          *discordgo.Session
	messages         []*discordgo.Message
	incomingMessages *chan discordgo.Message
	waitGroup        *sync.WaitGroup
	activated        bool
)

func init() {
	session = nil
	activated = false
	waitGroup = nil
	ch := make(chan discordgo.Message, 100000)
	incomingMessages = &ch
}

// SetToken creates a new Discord session and will automate some startup
// tasks if given enough information to do so.  Currently you can pass zero
// arguments and it will return an empty Discord session.
// There are 3 ways to call New:
//     With a single auth token - All requests will use the token blindly,
//         no verification of the token will be done and requests may fail.
//         IF THE TOKEN IS FOR A BOT, IT MUST BE PREFIXED WITH `BOT `
//         eg: `"Bot <token>"`
//     With an email and password - Discord will sign in with the provided
//         credentials.
//     With an email, password and auth token - Discord will verify the auth
//         token, if it is invalid it will sign in with the provided
//         credentials. This is the Discord recommended way to sign in.
//
// NOTE: While email/pass authentication is supported by DiscordGo it is
// HIGHLY DISCOURAGED by Discord. Please only use email/pass to obtain a token
// and then use that authentication token for all future connections.
// Also, doing any form of automation with a user (non Bot) account may result
// in that account being permanently banned from Discord.
func SetToken(token string) error {
	fmt.Printf("Token used: " + token)
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	session = s
	err = session.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}

	session.AddHandler(messageCreate)
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Printf("%v: %v", m.Author.Username, m.Content)
	if m.Content == "!beat! ping" {
		fmt.Println("Pong")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID != "155865120910737408" {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!beat! start" {
		fmt.Println("Starting scan")
		s.ChannelMessageSend(m.ChannelID, "Starting scan")
		getAllMessagesOfChannel(m.ChannelID, m.ID)
	}
}

func GetLatestMessages() chan discordgo.Message {

	if activated {
		old := *incomingMessages

		ch := make(chan discordgo.Message, 100000)
		incomingMessages = &ch
		close(old)
		return old
	} else {
		return nil
	}

}

func TurnOff() {
	fmt.Println("Closed Conection with discord")
	session.Close()
}

func getAllMessagesOfChannel(channelID string, lastMessageID string) {
	activated = true
	fmt.Println("Pasa aqui 1")
	if waitGroup == nil {
		fmt.Println("Pasa aqui 2")
		waitGroup = &sync.WaitGroup{}
	}

	messages, err := session.AllChannelMessages(channelID, lastMessageID)
	fmt.Println("Pasa aqui 3")
	var size int

	if err != nil {
		session.ChannelMessage(channelID, "failed, reason: "+err.Error())
		return
	} else {
		fmt.Println("Pasa aqui 4")
		size = len(messages)
		if size > 0 {
			lastValue := size - 1
			message := messages[lastValue]
			go getAllMessagesOfChannel(channelID, message.ID)
		}
	}

	go func() {
		fmt.Println("Pasa aqui 5")
		if size == 0 {
			fmt.Println("Pasa aqui 6")
			waitGroup.Wait()
			close(*incomingMessages)
			session.ChannelMessageSend(channelID, "---COMPLETED----")
			return
		}
		waitGroup.Add(1)
		print("Recieved chunk of data")
		for _, message := range messages {
			*incomingMessages <- *message
		}

		time.After(5 * time.Millisecond)
		waitGroup.Done()
	}()

}
