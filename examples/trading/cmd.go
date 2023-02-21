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

	fmt.Print("Start ...\n")
	err := client.Start(context.Background())
	defer client.Stop()
	if err != nil {
		log.Fatalf("Start fail: %s\n", err.Error())
		return
	}

	fmt.Print("Start success\n")

	prompt := `You:BOLL data changed: UpBand:[2.653 2.645 2.640 2.634 2.622 2.611 2.614 2.615 2.618 2.618 2.619 2.619 2.622 2.624 2.624 2.624 2.624 2.624 2.624 2.627], SMA:[2.605 2.603 2.601 2.599 2.596 2.594 2.595 2.596 2.598 2.599 2.600 2.599 2.599 2.598 2.598 2.598 2.598 2.598 2.598 2.599], DownBand:[2.557 2.561 2.562 2.564 2.570 2.577 2.575 2.577 2.579 2.579 2.581 2.580 2.575 2.572 2.572 2.572 2.572 2.572 2.571 2.571]
You:RSI data changed: [55.703 78.253 44.869 33.871 26.280 30.286 81.857 78.360 85.344 38.224 40.336 12.013 8.355 24.564 64.706 72.386 64.481 44.202 75.244 83.419]
You:There are currently no open positions
You:Analyze the data and generate only one trading command: /open_long_position, /open_short_position, /close_position or /no_action, the entity will execute the command and give you feedback.
AI:`
	fmt.Printf("%s", prompt)
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
