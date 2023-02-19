package main

import (
	"fmt"

	"github.com/Valgard/godotenv"
	"github.com/yubing744/chatgpt-go/pkg"
	"github.com/yubing744/chatgpt-go/pkg/config"
)

func main() {
	dotenv := godotenv.New()
	if err := dotenv.Load(".env.local"); err != nil {
		panic(err)
	}

	cfg := &config.Config{}
	client := pkg.NewChatgptClient(cfg)
	fmt.Printf("Hello World: %v", client)
}
