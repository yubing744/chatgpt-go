package pkg

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResponse(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := ``
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	msgs, err := client.parseResponse(resp)
	assert.NoError(t, err)
	assert.NotNil(t, msgs)
}

func TestParseResponseForUnmarshalError(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := `data: 2023-02-21 07:00:21.653311`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	msgs, err := client.parseResponse(resp)
	assert.NoError(t, err)
	assert.NotNil(t, msgs)
}

func TestParseResponseForDetail(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := `{"detail":"Too many requests in 1 hour. Try again later."}`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, "Too many requests in 1 hour. Try again later.", err.Error())
}

func TestParseResponseForServerError(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := `{"detail":{"message":"The server had an error while processing your request. Sorry about that! You can retry your request, or contact us through our help center at help.openai.com if the error persists. (Please include the request ID 985e0eeb2c44145e93637d2d79d416cf in your message.)","type":"server_error","param":null,"code":null}}`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, `{"detail":{"message":"The server had an error while processing your request. Sorry about that! You can retry your request, or contact us through our help center at help.openai.com if the error persists. (Please include the request ID 985e0eeb2c44145e93637d2d79d416cf in your message.)","type":"server_error","param":null,"code":null}}`, err.Error())
}

func TestParseResponseForInternalServerError(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := `Internal Server Error`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, "Internal Server Error", err.Error())
}

func TestParseResponseForLongMessagesShouldBeOk(t *testing.T) {
	client := NewChatgptClient("test", "test")

	body := `data: {"message": {"id": "d4cb6686-6ef7-46da-87d7-df16abbec928", "role": "assistant", "user": null, "create_time": null, "update_time": null, "content": {"content_type": "text", "parts": ["The term \"911\" is commonly associated with the emergency telephone number in the United States and Canada. When someone dials 911, it connects them with emergency services such as police, fire, or medical"]}, "end_turn": null, "weight": 1.0, "metadata": {"message_type": "next", "model_slug": "text-davinci-002-render-sha"}, "recipient": "all"}, "conversation_id": "45e1a523-c85a-4d11-96c5-b91e38f0ee83", "error": null}` + "\n"
	body = body + "data: [DONE]"

	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	msgs, err := client.parseResponse(resp)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, `The term "911" is commonly associated with the emergency telephone number in the United States and Canada. When someone dials 911, it connects them with emergency services such as police, fire, or medical`, msgs[0].Text)
}
