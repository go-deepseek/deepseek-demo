package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-deepseek/deepseek"
)

const DEEPSEEK_API_KEY = `sk-123cd456c78d9be0b123de45cf6789b0` // replace with valid one

func main() {
	client := deepseek.NewClientWithTimeout(DEEPSEEK_API_KEY, 120)

	fmt.Println("This is demo for deepseek.  Type `bye` to exit.")

	for {
		fmt.Print(">>> ")

		reader := bufio.NewReader(os.Stdin)
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		inMsg := string(lineBytes)
		fmt.Println(inMsg)

		if strings.ToLower(inMsg) == "bye" || strings.ToLower(inMsg) == "exit" || strings.ToLower(inMsg) == "quit" {
			fmt.Println("bye!!!")
			os.Exit(0)
		}

		// req := &deepseek.DeepseekChatRequest{
		// 	Model: "deepseek-chat",
		// 	Messages: []*deepseek.Message{
		// 		{
		// 			Role:    "system",
		// 			Content: "You are a helpful assistant.",
		// 		},
		// 		{
		// 			Role:    "user",
		// 			Content: inMsg,
		// 		},
		// 	},
		// 	Stream: true,
		// }

		req := &deepseek.DeepseekChatRequest{
			Model: "deepseek-reasoner",
			Messages: []*deepseek.Message{
				{
					Role:    "user",
					Content: inMsg,
				},
			},
			Stream: true,
		}

		iter, err := client.StreamChatCompletionsReasoner(req)
		if err != nil {
			panic(err)
		}

		// fmt.Println(iter.Choices[0].Message.ReasoningContent)
		// fmt.Println()
		// fmt.Println()
		// fmt.Println(iter.Choices[0].Message.Content)

		for {
			msg := iter.Next()
			if msg == nil {
				break
			}
			fmt.Print(msg.Choices[0].Delta.ReasoningContent)
			fmt.Print(msg.Choices[0].Delta.Content)
		}
		fmt.Println()
	}

}
