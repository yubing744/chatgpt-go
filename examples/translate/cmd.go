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
		Timeout:  time.Second * 300,
		Debug:    false,
	}
	client := pkg.NewChatgptClient(cfg)

	fmt.Print("Starting ...\n")
	err := client.Start(context.Background())
	defer client.Stop()
	if err != nil {
		log.Fatalf("Start fail: %s\n", err.Error())
		return
	}

	fmt.Print("Start success\n")

	prompt := "翻译成英文：你还需要哪些指标帮助决策交易命令？"
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
