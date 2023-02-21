package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func (client *ChatgptClient) Ask(ctx context.Context, prompt string, conversationId *string, parentId *string) (*AskResult, error) {
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

		msgs, err := client.parseResponse(resp)
		if err != nil {
			return nil, err
		}

		if len(msgs) > 0 {
			result.Data = msgs[len(msgs)-1]
		}

		return result, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return nil, errors.Errorf("Error in ask: %s", string(body))
}

func (client *ChatgptClient) parseResponse(response *http.Response) ([]*Message, error) {
	log := client.logger

	messages := make([]*Message, 0)

	log.Printf("\n")
	log.Printf("Parse response ")

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()

		log.Printf(".")

		if client.debug {
			log.Printf("new line: %s\n", line)
		}

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			log.Printf("line: %s\n", line)

			line = strings.ReplaceAll(line, `\"`, `"`)
			line = strings.ReplaceAll(line, `\'`, `'`)
			line = strings.ReplaceAll(line, `\\`, `\`)

			var data struct {
				Detail string `json:"detail"`
			}
			err := json.Unmarshal([]byte(line), &data)
			if err != nil {
				return nil, errors.New(line)
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
			log.Printf("Error in Unmarshal: %s\n", line)
			continue
		}

		if !checkFields(parsedLine) {
			log.Printf("Field missing\n")
			log.Printf("%v", parsedLine)
			continue
		}

		messageContextType := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["content_type"].(string)
		if messageContextType == "text" {
			message := parsedLine["message"].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0]
			conversationID := parsedLine["conversation_id"].(string)
			parentID := parsedLine["message"].(map[string]interface{})["id"].(string)
			messages = append(messages, &Message{
				ConversationID: conversationID,
				ParentID:       parentID,
				Text:           fmt.Sprintf("%v", message),
			})
		} else {
			log.Printf("not support message type: %s\n", messageContextType)
		}
	}

	log.Printf("\n")

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
