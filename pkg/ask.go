package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Message struct {
	Text           string
	ConversationID string
	ParentID       string
}

type AskResult struct {
	Code   int
	Detail string `json:"detail"`
	Data   *Message
}

func (client *ChatgptClient) Ask(ctx context.Context, prompt string, conversationId *string, parentId *string, timeout time.Duration) (*AskResult, error) {
	url := fmt.Sprintf("%s/%s", client.baseURL, "api/conversation")
	headers := http.Header{
		"Accept": {"application/json; charset=utf-8"},
	}

	data := map[string]interface{}{
		"action": "next",
		"messages": []map[string]interface{}{
			{
				"id":   uuid.New().String(),
				"role": "user",
				"content": map[string]interface{}{
					"content_type": "text",
					"parts":        []string{prompt},
				},
			},
		},
		"model": "text-davinci-002-render-sha",
	}

	if conversationId != nil {
		data["conversation_id"] = *conversationId
	}

	if parentId != nil {
		data["parent_message_id"] = *parentId
	} else {
		data["parent_message_id"] = uuid.New().String()
	}

	payload, _ := json.Marshal(data)
	resp, err := client.session.Post(url, headers, payload, true)
	if err != nil {
		return nil, errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		result := &AskResult{
			Code:   0,
			Detail: "",
		}

		msgs, err := parseResponse(resp)
		if err != nil {
			result.Code = 1
			result.Detail = err.Error()
			return result, nil
		}

		if len(msgs) > 0 {
			result.Data = msgs[len(msgs)-1]
		}

		return result, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return nil, errors.Errorf("Error in ask: %s", string(body))
}

func parseResponse(response *http.Response) ([]*Message, error) {
	messages := make([]*Message, 0)

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("raw line: %s\n", line)

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			var data struct {
				Detail string `json:"detail"`
			}
			err := json.Unmarshal([]byte(line), &data)
			if err != nil {
				return nil, err
			}

			return nil, errors.New(data.Detail)
		}

		line = strings.TrimPrefix(line, "data: ")
		if line == "[DONE]" {
			break
		}

		line = strings.ReplaceAll(line, `\"`, `"`)
		line = strings.ReplaceAll(line, `\'`, `'`)
		line = strings.ReplaceAll(line, `\\`, `\`)

		var parsedLine map[string]interface{}
		err := json.Unmarshal([]byte(line), &parsedLine)
		if err != nil {
			continue
		}
		if !checkFields(parsedLine) {
			fmt.Println("Field missing")
			fmt.Println(parsedLine)
			continue
		}

		messageContextType := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["content_type"].(string)
		if messageContextType == "test" {
			message := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0]
			conversationID := parsedLine["conversation_id"].(string)
			parentID := parsedLine["message"].(map[string]interface{})["id"].(string)
			messages = append(messages, &Message{
				ConversationID: conversationID,
				ParentID:       parentID,
				Text:           fmt.Sprintf("%v", message),
			})
		} else {
			fmt.Printf("not support message type: %s", messageContextType)
		}
	}

	return messages, nil
}

func checkFields(parsedLine map[string]interface{}) bool {
	_, messageExists := parsedLine["message"]
	_, conversationIDExists := parsedLine["conversation_id"]
	_, messageContentExists := parsedLine["message"].(map[string]interface{})["content"]
	_, messageContentTypeExists := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["content_type"]
	_, messagePartsExists := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["parts"]
	if messageExists && conversationIDExists && messageContentExists && messageContentTypeExists && messagePartsExists {
		return true
	}
	return false
}
