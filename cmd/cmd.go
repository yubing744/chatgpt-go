package main

import (
	"fmt"

	"github.com/yubing744/chatgpt-go/pkg"
)

func main() {
	client := pkg.NewChatgptClient()
	fmt.Printf("Hello World: %v", client)
}
