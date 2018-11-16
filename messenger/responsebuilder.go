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

func QRText(title string, payload string) QuickReply {
	return QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}
