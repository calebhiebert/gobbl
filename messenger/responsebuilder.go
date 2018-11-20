/*
	responsebuilder.go

	This file contains utilities for crafting Messenger response objects
*/

package fb

import (
	"errors"
	"math/rand"
	"time"

	"github.com/calebhiebert/gobbl"
)

type MBResponse struct {
	Messages      []OutgoingMessage
	QuickReplies  []QuickReply
	MinTypingTime []time.Duration
}

type MBImmediateResponse struct {
	MBResponse
	Context     *gbl.Context
	Integration *MessengerIntegration
}

// CreateResponse will return a pre-populated messenger response object and add it to the context
func CreateResponse(c *gbl.Context) *MBResponse {
	existingResponse := c.R.(*MBResponse)

	if existingResponse.Messages == nil {
		existingResponse.Messages = []OutgoingMessage{}
	}

	if existingResponse.QuickReplies == nil {
		existingResponse.QuickReplies = []QuickReply{}
	}

	if existingResponse.MinTypingTime == nil {
		existingResponse.MinTypingTime = []time.Duration{time.Second}
	}

	return existingResponse
}

// CreateImmediateResponse will create a response object that can be sent manually at any time
func CreateImmediateResponse(c *gbl.Context) *MBImmediateResponse {
	integration := c.Integration.(*MessengerIntegration)

	immediateResponse := MBImmediateResponse{
		Integration: integration,
		Context:     c,
		MBResponse: MBResponse{
			Messages:      []OutgoingMessage{},
			QuickReplies:  []QuickReply{},
			MinTypingTime: []time.Duration{time.Second},
		},
	}

	return &immediateResponse
}

func (im *MBImmediateResponse) Send() error {
	if im.Context.User.ID == "" {
		return errors.New("Missing user ID!")
	}

	return im.Integration.doResponse(im.Context.User.ID, &im.MBResponse)
}

// M adds a new message to the response and returns it
func (m *MBResponse) M(om *OutgoingMessage) {
	m.Messages = append(m.Messages, *om)
}

func (m *MBResponse) Text(text string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Text: text,
	})
}

func (m *MBResponse) RandomText(text ...string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Text: text[rand.Intn(len(text))],
	})
}

func (m *MBResponse) TypingTime(mtt ...time.Duration) {
	m.MinTypingTime = mtt
}

func (m *MBResponse) Image(url string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Attachment: &OutgoingAttachment{
			Type: "image",
			Payload: TemplatePayload{
				URL:        url,
				IsReusable: true,
			},
		},
	})
}

func (m *MBResponse) QR(qr ...QuickReply) {
	for _, quickReply := range qr {
		m.QuickReplies = append(m.QuickReplies, quickReply)
	}
}

func QRText(title string, payload string) QuickReply {
	return QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}
