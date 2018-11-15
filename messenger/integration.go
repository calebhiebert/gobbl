/*
	integration.go

	This file contains the integration for Facebook Messenger
*/

// Package fb impliments a -ho-hook integration for Facebook Messenger
package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	".."
)

type MessengerIntegration struct {
	API         *MessengerAPI
	Bot         *gbl.Bot
	VerifyToken string
}

// GenericRequest extracts a generic request from a facebook webhook request.
// The generic response Text property is set on the following rules
// If the event is a message, it will the the message text.
// If the event is a quick reply, it will be the quick reply payload.
// If the event is a postback, it will first try to use the postback payload, but will fall back to the payload title.
// If the event is a referral, it will be the ref property.
// This method will also sest the fb:eventtype flag on the context, it will be one of the following values:
// quickreply, message, payload, referral
func (m *MessengerIntegration) GenericRequest(c *gbl.Context) (gbl.GenericRequest, error) {
	genericRequest := gbl.GenericRequest{}
	fbRequest := c.RawRequest.(MessagingItem)

	// Check for a message id
	if fbRequest.Message.MID != "" {

		// Check for a quickreply payload
		if fbRequest.Message.QuickReply.Payload != "" {
			genericRequest.Text = fbRequest.Message.QuickReply.Payload
			c.Flag("fb:eventtype", "quickreply")
		} else {
			genericRequest.Text = fbRequest.Message.Text
			c.Flag("fb:eventtype", "message")
		}

		// Check for a postback title
	} else if fbRequest.Postback.Title != "" {

		// First try to use the postback payload
		if fbRequest.Postback.Payload != "" {
			genericRequest.Text = fbRequest.Postback.Payload

			// Fall back to the postback title on issues
		} else {
			genericRequest.Text = fbRequest.Postback.Title
		}
		c.Flag("fb:eventtype", "payload")

		// Check for a referral
	} else if fbRequest.Referral.Ref != "" {
		genericRequest.Text = fbRequest.Referral.Ref
		c.Flag("fb:eventtype", "referral")
	}

	return genericRequest, nil
}

// User will extract a user's psid from a facebook webhook request
// It will check if the message is an echo, to make sure the correct id is always selected
func (m *MessengerIntegration) User(c *gbl.Context) (gbl.User, error) {
	user := gbl.User{}

	fbRequest := c.RawRequest.(MessagingItem)

	if fbRequest.Message.MID != "" {
		if fbRequest.Message.IsEcho {
			user.ID = fbRequest.Recipient.ID
		} else {
			user.ID = fbRequest.Sender.ID
		}
	} else if fbRequest.Postback.Title != "" {
		user.ID = fbRequest.Sender.ID
	} else if fbRequest.Referral.Ref != "" {
		user.ID = fbRequest.Sender.ID
	} else {
		return user, errors.New("Unable to determine facebook event type")
	}

	return user, nil
}

// Respond will respond to messages using a MBResponse struct
// Optionally it will set the typing icon and wait for the duration set in the MBResponse.MinTypingTime slice
// If the MinTypingTime slice contains only one value, that value will be used for all messages
// If the MinTypingTime slice contains a number of values equal to the number of messages, one value per message will be used
// If the MinTypingTime slice contains a number of values that is not 0, 1, or the number of messages, an error will be returned
// The default typing time is set to 1 second
func (m *MessengerIntegration) Respond(c *gbl.Context) (*interface{}, error) {
	if c.User.ID == "" {
		return nil, errors.New("Unable to respond, user id missing")
	}

	response := c.R.(*MBResponse)

	// Append quick replies to the last message if they exist
	if len(response.QuickReplies) > 0 {
		response.Messages[len(response.Messages)-1].QuickReplies = response.QuickReplies
	}

	// Loop through each message and send it
	for idx, msg := range response.Messages {

		// Check to see if the typing indicator should be set
		if len(response.MinTypingTime) == 1 || len(response.MinTypingTime) == len(response.Messages) {
			_, err := m.API.SenderAction(&User{
				ID: c.User.ID,
			}, SenderActionTypingOn)
			if err != nil {
				fmt.Printf("Error while setting typing %+v\n", err)
			}

			// Sleep for the appropriate amount of time
			if len(response.MinTypingTime) == 1 {
				time.Sleep(response.MinTypingTime[0])
			} else if len(response.MinTypingTime) == len(response.Messages) {
				time.Sleep(response.MinTypingTime[idx])
			}
		} else if len(response.MinTypingTime) != 0 && len(response.MinTypingTime) != len(response.Messages) {
			return nil, errors.New("Typing time mismatch")
		}

		// Send the message
		_, err := m.API.SendMessage(&User{
			ID: c.User.ID,
		}, MessageTypeResponse, &msg)
		if err != nil {
			fmt.Printf("Error while sending message %+v\n", err)
		}
	}

	return nil, nil
}

// ServeHTTP is a http request handler that is specifically built for accepting facebook webhook requests
func (m *MessengerIntegration) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Recovery function
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
			rw.WriteHeader(http.StatusInternalServerError)

			jsonErr, err := json.Marshal(map[string]interface{}{
				"error": r,
			})
			if err != nil {
				rw.Write([]byte(fmt.Sprintf("%+v", r)))
				return
			}

			rw.Write(jsonErr)
			return
		}
	}()

	// Check the request method so webhook verification can be completed
	if req.Method == "GET" {
		mode := req.URL.Query()["hub.mode"][0]
		token := req.URL.Query()["hub.verify_token"][0]
		challenge := req.URL.Query()["hub.challenge"][0]

		if mode == "subscribe" && token == m.VerifyToken {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(challenge))
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// We are receiving a webhook request (probably)
	} else if req.Method == "POST" {
		defer req.Body.Close()

		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Error reading request %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var webhookRequest WebhookRequest

		err = json.Unmarshal(requestBody, &webhookRequest)
		if err != nil {
			fmt.Printf("Error parsing json request %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages := []MessagingItem{}

		// Loop through facebook's request messages
		for _, entry := range webhookRequest.Entry {
			for _, msg := range entry.Messaging {
				messages = append(messages, msg)
			}
		}

		// TODO execute all received messages in parallel with goroutines
		// Execute each message with the bot
		for _, message := range messages {
			inputCtx := gbl.InputContext{
				RawRequest:  message,
				Integration: m,
				Response: &MBResponse{
					Messages:      []OutgoingMessage{},
					QuickReplies:  []QuickReply{},
					MinTypingTime: []time.Duration{time.Second},
				},
			}

			m.Bot.Execute(&inputCtx)
		}

		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Listen will start a server listening for facebook webhook requests
func (m *MessengerIntegration) Listen(server *http.Server, bot *gbl.Bot) {

	server.Handler = m
	m.Bot = bot

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
