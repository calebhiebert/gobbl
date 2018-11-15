/*
	responsebuilder.go

	This file contains utilities for crafting Messenger response objects
*/

package fb

import (
	"time"

	".."
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
