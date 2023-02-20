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
		Debug:    true,
	}
	client := pkg.NewChatgptClient(cfg)

	fmt.Print("Login ...\n")
	err := client.Login()
	if err != nil {
		log.Fatalf("Login fail: %s\n", err.Error())
		return
	}

	fmt.Print("Login success\n")

	prompt := "openAI API 接口 模型温度如何设置？"
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
