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
	"runtime/debug"
	"time"

	"github.com/calebhiebert/gobbl"
	. "github.com/logrusorgru/aurora"
)

type MessengerIntegration struct {
	API            *MessengerAPI
	Bot            *gbl.Bot
	VerifyToken    string
	DevMode        bool
	DisableTyping  bool
	EnableRecovery bool
	Always200      bool
}

// RetargetRawRequest is the request type that will be set
// as the raw request when retargeting
type RetargetRawRequest struct {
	Type string
	PSID string
	Args interface{}
}

// GenericRequest extracts a generic request from a facebook webhook request.
// The generic response Text property is set on the following rules
// If the event is a message, it will the the message text.
// If the event is a quick reply, it will be the quick reply payload.
// If the event is a postback, it will first try to use the postback payload, but will fall back to the payload title.
// If the event is a referral, it will be the ref property.
// This method will also set the fb:eventtype flag on the context, it will be one of the following values:
// quickreply, message, payload, referral, coordinates, attachment
func (m *MessengerIntegration) GenericRequest(c *gbl.Context) (gbl.GenericRequest, error) {
	genericRequest := gbl.GenericRequest{}

	switch c.RawRequest.(type) {
	case *MessagingItem:
		fbRequest := c.RawRequest.(*MessagingItem)

		if fbRequest.IsStandby {
			c.Trace("CURRENT MESSAGE IS STANDBY")
			c.Flag("fb:isstandby", true)
		}

		// Check for a message id
		if fbRequest.Message.MID != "" {

			if fbRequest.Message.IsEcho == true {
				c.Trace("CURRENT MESSAGE IS ECHO")
				c.Flag("fb:isecho", true)
			}

			if fbRequest.Message.AppID != 0 {
				c.Flag("fb:sendingappid", fmt.Sprintf("%d", fbRequest.Message.AppID))
			}

			// Check for a quickreply payload
			if fbRequest.Message.QuickReply.Payload != "" {
				genericRequest.Text = fbRequest.Message.QuickReply.Payload
				c.Flag("fb:eventtype", "quickreply")

				// Check for message text
			} else if fbRequest.Message.Text != "" {
				genericRequest.Text = fbRequest.Message.Text
				c.Flag("fb:eventtype", "message")

				// Check for coordinates
			} else if len(fbRequest.Message.Attachments) > 0 && fbRequest.Message.Attachments[0].Payload.Coordinates.Lat != 0 {
				genericRequest.Text = fmt.Sprintf("COORDS LAT %f LONG %f",
					fbRequest.Message.Attachments[0].Payload.Coordinates.Lat,
					fbRequest.Message.Attachments[0].Payload.Coordinates.Long)
				c.Flag("fb:eventtype", "coordinates")
				c.Flag("fb:location", fbRequest.Message.Attachments[0].Payload.Coordinates)

				// Check for an attachment
			} else if len(fbRequest.Message.Attachments) > 0 && fbRequest.Message.Attachments[0].Payload.URL != "" {
				genericRequest.Text = fbRequest.Message.Attachments[0].Payload.URL
				c.Flag("fb:eventtype", "attachment")
				c.Flag("fb:attachmenturl", fbRequest.Message.Attachments[0].Payload.URL)
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

			if fbRequest.Postback.Referral.Ref != "" {
				c.Flag("fb:referral", fbRequest.Postback.Referral)
			}

			// Check for a referral
		} else if fbRequest.Referral.Ref != "" {
			genericRequest.Text = fbRequest.Referral.Ref
			c.Flag("fb:eventtype", "referral")
			c.Flag("fb:referral", fbRequest.Referral)
		} else if fbRequest.TakeThreadControl.PreviousOwnerAppID != 0 {
			c.Flag("fb:eventtype", "take_thread_control")
			c.Flag("fb:handover:metadata", fbRequest.TakeThreadControl.Metadata)
		}
	case RetargetRawRequest:
		genericRequest.Text = "RETARGET"
		c.Flag("fb:eventtype", "retarget")
	}

	return genericRequest, nil
}

// User will extract a user's psid from a facebook webhook request
// It will check if the message is an echo, to make sure the correct id is always selected
func (m *MessengerIntegration) User(c *gbl.Context) (gbl.User, error) {
	user := gbl.User{}

	switch c.RawRequest.(type) {
	case *MessagingItem:
		fbRequest := c.RawRequest.(*MessagingItem)

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
		} else if fbRequest.TakeThreadControl.PreviousOwnerAppID != 0 {
			user.ID = fbRequest.Sender.ID
		} else {
			return user, errors.New("Unable to determine facebook event type")
		}
	case RetargetRawRequest:
		user.ID = c.RawRequest.(RetargetRawRequest).PSID
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

	m.doResponse(c.User.ID, response, c)

	return nil, nil
}

// ProcessMessage will execute a single facebook message in the bot
func (m *MessengerIntegration) ProcessMessage(message *MessagingItem, errChan chan error) {

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(Red("FB PANIC: "), Red(r))
			fmt.Println(Red(string(debug.Stack())))

			errChan <- fmt.Errorf("%v", r)
		}
	}()

	// Request processing
	inputCtx := gbl.InputContext{
		RawRequest:  message,
		Integration: m,
		Response:    &MBResponse{},
	}

	m.Bot.Execute(&inputCtx)

	errChan <- nil
}

// ProcessWebhookRequest will process an incoming facebook webhook request
// each messaging item in the request will be processed in it's own goroutine
func (m *MessengerIntegration) ProcessWebhookRequest(request *WebhookRequest) []error {
	messages := []MessagingItem{}

	// Loop through facebook's request messages
	for _, entry := range request.Entry {
		for _, msg := range entry.Messaging {
			messages = append(messages, msg)
		}

		for _, sby := range entry.Standby {
			sby.IsStandby = true
			messages = append(messages, sby)
		}
	}

	errChan := make(chan error)

	// Execute each message with the bot
	for _, message := range messages {
		go m.ProcessMessage(&message, errChan)
	}

	var executionErrs = []error{}

	for i := 0; i < len(messages); i++ {
		err := <-errChan
		if err != nil {
			executionErrs = append(executionErrs, err)
		}
	}

	return executionErrs
}

// ServeHTTP is a http request handler that is specifically built for accepting facebook webhook requests
func (m *MessengerIntegration) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.EnableRecovery || m.Always200 {
		// Recovery function
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(Red("FB PANIC: "), Red(r))
				fmt.Println(Red(string(debug.Stack())))

				if m.Always200 {
					rw.WriteHeader(http.StatusOK)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
				}

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
	}

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
			if m.Always200 {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		var webhookRequest WebhookRequest

		err = json.Unmarshal(requestBody, &webhookRequest)
		if err != nil {
			fmt.Printf("Error parsing json request %+v", err)
			if m.Always200 {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		_ = m.ProcessWebhookRequest(&webhookRequest)

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

// Retarget will send a message through the bot with the text RETARGET
// this can be used to route through the bot in a normal way while
// sending retargeting messages
func (m *MessengerIntegration) Retarget(psid, retargetType string, args interface{}) error {
	inputCtx := gbl.InputContext{
		RawRequest: RetargetRawRequest{
			Type: retargetType,
			PSID: psid,
			Args: args,
		},
		Integration: m,
		Response:    &MBResponse{},
	}

	_, err := m.Bot.Execute(&inputCtx)

	return err
}

func (m *MessengerIntegration) doResponse(psid string, response *MBResponse, c *gbl.Context) error {
	// Append quick replies to the last message if they exist
	if len(response.QuickReplies) > 0 && len(response.Messages) > 0 {
		response.Messages[len(response.Messages)-1].QuickReplies = response.QuickReplies
	}

	// Loop through each message and send it
	for idx, msg := range response.Messages {

		// Check to see if the typing indicator should be set
		if !m.DisableTyping && len(response.MinTypingTime) == 1 || len(response.MinTypingTime) == len(response.Messages) {
			_, err := m.API.SenderAction(&User{
				ID: psid,
			}, SenderActionTypingOn)
			if err != nil {
				c.Log(10, fmt.Sprintf("Error setting typing %v", err), "FBResponse")
				return err
			}
			c.Log(50, "Set typing success", "FBResponse")

			// Sleep for the appropriate amount of time
			if len(response.MinTypingTime) == 1 {
				time.Sleep(response.MinTypingTime[0])
			} else if len(response.MinTypingTime) == len(response.Messages) {
				time.Sleep(response.MinTypingTime[idx])
			}
		} else if !m.DisableTyping && len(response.MinTypingTime) != 0 && len(response.MinTypingTime) != len(response.Messages) {
			return errors.New("Typing time mismatch")
		}

		// Send the message
		resp, err := m.API.SendMessage(&User{
			ID: psid,
		}, MessageTypeResponse, &msg)
		if err != nil {
			c.Log(10, fmt.Sprintf("Error sending message %v", err), "FBResponse")
			if m.DevMode {
				m.API.SendMessage(&User{
					ID: psid,
				}, MessageTypeResponse, &OutgoingMessage{
					Text: fmt.Sprintf("Message Error!\n%+v", err),
				})
			}

			return err
		}

		c.Log(50, fmt.Sprintf("Sent message %v", resp), "FBResponse")
	}

	return nil
}
