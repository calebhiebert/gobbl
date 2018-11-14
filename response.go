package main

type Response interface {
	Text(t string) Response
}

type BasicResponse struct {
	messages []string
}

func (br BasicResponse) Text(t string) Response {
	br.messages = append(br.messages, t)
	return br
}
