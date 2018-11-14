package fb

type ResponseBuilder struct {
	messages     []RBMessage
	quickReplies []QuickReply
}

func (rb *ResponseBuilder) Message(m RBMessage) *ResponseBuilder {
	rb.messages = append(rb.messages, m)
	return rb
}

func (rb *ResponseBuilder) QuickReply(qr QuickReply) *ResponseBuilder {
	rb.quickReplies = append(rb.quickReplies, qr)
	return rb
}

func (rb *ResponseBuilder) M() *RBMessage {
	m := RBMessage{}

	rb.Message(m)

	return &m
}

func (rb *ResponseBuilder) Build() *[]OutgoingMessage {
	messages := []OutgoingMessage{}

	for _, msg := range rb.messages {
		messages = append(messages, *msg.Build())
	}

	if len(rb.quickReplies) > 0 {
		messages[len(messages)-1].QuickReplies = rb.quickReplies
	}

	return &messages
}

type RBMessage struct {
	message        OutgoingMessage
	attachmentType string
}

func (m *RBMessage) Text(text string) *RBMessage {
	m.message.Text = text
	return m
}

func (m *RBMessage) Build() *OutgoingMessage {
	return &m.message
}
