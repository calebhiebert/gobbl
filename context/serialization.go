package bctx

import "encoding/json"

func encodeContext(ctx *BotContext) (string, error) {
	jsonBytes, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func decodeContext(jsonString string) (BotContext, error) {
	var bctx BotContext

	err := json.Unmarshal([]byte(jsonString), &bctx)
	if err != nil {
		return bctx, err
	}

	return bctx, nil
}
