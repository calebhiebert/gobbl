/*
	responsebuilder.go

	This file contains utilities for crafting Messenger response objects
*/

package fb

import (
	"time"

	"github.com/calebhiebert/gobbl"
)

type MBResponse struct {
	Messages      []OutgoingMessage
	QuickReplies  []QuickReply
	MinTypingTime []time.Duration
}

// CreateResponse will return a pre-populated messenger response object and add it to the context
func CreateResponse(c *gbl.Context) *MBResponse {
	r := &MBResponse{
		Messages:      []OutgoingMessage{},
		QuickReplies:  []QuickReply{},
		MinTypingTime: []time.Duration{time.Second},
	}

	c.R = r

	return r
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

func (m *MBResponse) QR(qr QuickReply) {
	m.QuickReplies = append(m.QuickReplies, qr)
}

func QRText(title string, payload string) QuickReply {
	return QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}
