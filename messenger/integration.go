package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	".."
)

type MessengerIntegration struct {
	API         *MessengerAPI
	Bot         *cpn.Bot
	VerifyToken string
}

func (m *MessengerIntegration) GenericRequest(c *cpn.Context) (*cpn.GenericRequest, error) {

	genericRequest := cpn.GenericRequest{}

	fbRequest := (*c.RawRequest).(MessagingItem)

	if fbRequest.Message.MID != "" {
		if fbRequest.Message.QuickReply.Payload != "" {
			genericRequest.Text = fbRequest.Message.QuickReply.Payload
			c.Flag("fb:eventtype", "quickreply")
		} else {
			genericRequest.Text = fbRequest.Message.Text
			c.Flag("fb:eventtype", "message")
		}
	} else if fbRequest.Postback.Title != "" {
		if fbRequest.Postback.Payload != "" {
			genericRequest.Text = fbRequest.Postback.Payload
		} else {
			genericRequest.Text = fbRequest.Postback.Title
		}
		c.Flag("fb:eventtype", "payload")
	} else if fbRequest.Referral.Ref != "" {
		genericRequest.Text = fbRequest.Referral.Ref
		c.Flag("fb:eventtype", "referral")
	}

	return &genericRequest, nil
}

func (m *MessengerIntegration) User(c *cpn.Context) (*cpn.User, error) {
	user := cpn.User{}

	fbRequest := (*c.RawRequest).(MessagingItem)

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
		return nil, errors.New("Unable to determine facebook event type")
	}

	return &user, nil
}

func (m *MessengerIntegration) Respond(c *cpn.Context) (*interface{}, error) {
	if c.User == nil {
		return nil, errors.New("Unable to respond, user missing")
	}

	if c.User.ID == "" {
		return nil, errors.New("Unable to respond, user id missing")
	}

	responseBuilder := c.R.(ResponseBuilder)

	messages := responseBuilder.Build()

	fmt.Printf("Response Builder Messages %d", len(*messages))

	for _, msg := range *messages {
		result, err := m.API.SendMessage(&User{
			ID: c.User.ID,
		}, MessageTypeResponse, &msg)
		if err != nil {
			fmt.Printf("Error while sending message %+v\n", err)
		}

		fmt.Printf("Sent Message %+v\n", result)
	}

	return nil, nil
}

func (m *MessengerIntegration) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		mode := req.URL.Query()["hub.mode"][0]
		token := req.URL.Query()["hub.verify_token"][0]
		challenge := req.URL.Query()["hub.challenge"][0]

		if mode == "subscribe" && token == m.VerifyToken {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(challenge))
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
		}
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

		for _, entry := range webhookRequest.Entry {
			for _, msg := range entry.Messaging {
				messages = append(messages, msg)
			}
		}

		for _, message := range messages {
			inputCtx := cpn.InputContext{
				RawRequest:  message,
				Integration: m,
				Response:    ResponseBuilder{},
			}

			m.Bot.Execute(&inputCtx)
		}

		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (m *MessengerIntegration) Listen(server *http.Server, bot *cpn.Bot) {

	server.Handler = m
	m.Bot = bot

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
