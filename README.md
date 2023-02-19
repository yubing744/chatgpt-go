# chatgpt-go
chatGPT golang client translated from https://github.com/acheong08/ChatGPT


## Usage

Config .env.local file
```
CHATGPT_EMAIL="your chat gpt account"
CHATGPT_PASSWORD="your chat gpt password"
```

Run
``` bash
make run
```


## Example

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
	"github.com/yubing744/chatgpt-go/pkg/config"
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

	cfg := &config.Config{
		Email:    email,
		Password: password,
		Proxy:    "",
	}
	client := pkg.NewChatgptClient(cfg)

	fmt.Print("Login ...\n")
	err := client.Login()
	if err != nil {
		log.Fatalf("Login fail: %s\n", err.Error())
		return
	}

	fmt.Print("Login success\n")

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