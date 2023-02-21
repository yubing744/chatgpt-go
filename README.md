# chatgpt-go
chatGPT golang client translated from https://github.com/acheong08/ChatGPT

## Installation

```shell script
go get github.com/yubing744/chatgpt-go
```

## Usage

Config .env.local file
```
CHATGPT_EMAIL="your chat gpt account"
CHATGPT_PASSWORD="your chat gpt password"
```

``` go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Valgard/godotenv"
	"github.com/yubing744/chatgpt-go/pkg"
)

func main() {
	dotenv := godotenv.New()
	if err := dotenv.Load(".env.local"); err != nil {
		panic(err)
	}

	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	if email == "" || password == "" {
		log.Panic("CHATGPT_EMAIL or CHATGPT_PASSWORD not set in .env.local")
	}

	client := pkg.NewChatgptClient(email, password)

	fmt.Print("Starting ...\n")
	err := client.Start(context.Background())
	defer client.Stop()

	if err != nil {
		log.Fatalf("Start fail: %s\n", err.Error())
		return
	}

	fmt.Print("Start success\n")

	prompt := "Hello"
	fmt.Printf("You: %s", prompt)
	result, err := client.Ask(context.Background(), prompt, nil, nil, time.Second*5)
	if err != nil {
		fmt.Printf("Ask fail: %s\n", err.Error())
		return
	}

	if result.Code == 0 {
		fmt.Printf("AI: %s\n", result.Data.Text)
	}

	fmt.Print("Done\n")
}
```