package cap

type BasicResponse struct {
	messages []string
}

func (br *BasicResponse) Text(t string) *BasicResponse {
	br.messages = append(br.messages, t)
	return br
}
