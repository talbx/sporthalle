package fn

import (
	"context"
	"fmt"
	event2 "github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_receivePositive(t *testing.T) {
	// given
	data := "sporthalle"

	// when
	resp := sendMessageWithData(data)

	// then
	assert.Error(t, resp)
	assert.Equal(t, "test scenario ends here", resp.Error())
}

func Test_receiveNegative(t *testing.T) {
	// given
	data := "idc"

	// when
	resp := sendMessageWithData(data)

	// then
	assert.Error(t, resp)
	assert.Equal(t, fmt.Sprintf("could not do anything with msg %s", data), resp.Error())
}

func sendMessageWithData(data string) error {
	m := PubSubMessage{
		Data: []byte(data),
	}
	msg := MessagePublishedData{
		Message: m,
	}

	e := event2.New()
	e.SetDataContentType("application/json")
	e.SetData(e.DataContentType(), msg)

	ctx := context.WithValue(context.Background(), "test", "yes")
	resp := receive(ctx, e)
	return resp
}
