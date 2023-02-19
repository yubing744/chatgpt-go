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
	Data   []*Message
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

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		if strings.HasPrefix(string(body), "data:") {
			result := AskResult{
				Code:   0,
				Detail: "",
			}

			result.Data = parseResponse(resp)

			return &result, nil
		} else {
			result := AskResult{
				Code: 1,
				Data: nil,
			}

			err := json.Unmarshal(body, &result)
			if err != nil {
				return nil, errors.Wrap(err, "Error in Unmarshal body")
			}

			return &result, nil
		}
	}

	return nil, errors.New(string(body))
}

func parseResponse(response *http.Response) []*Message {
	messages := make([]*Message, 0)

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()[2 : len(scanner.Text())-1]
		if line == "" {
			continue
		}
		if strings.Contains(line, "data: ") {
			line = line[6:]
		}
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
		message := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0]
		conversationID := parsedLine["conversation_id"].(string)
		parentID := parsedLine["message"].(map[string]interface{})["id"].(string)
		messages = append(messages, &Message{
			ConversationID: conversationID,
			ParentID:       parentID,
			Text:           fmt.Sprintf("%v", message),
		})
	}

	return messages
}

func checkFields(parsedLine map[string]interface{}) bool {
	_, messageExists := parsedLine["message"]
	_, conversationIDExists := parsedLine["conversation_id"]
	_, messageContentExists := parsedLine["message"].(map[string]interface{})["content"]
	_, messagePartsExists := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["parts"]
	if messageExists && conversationIDExists && messageContentExists && messagePartsExists {
		return true
	}
	return false
}
