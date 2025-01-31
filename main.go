package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
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

		if strings.ToLower(inMsg) == "bye" || strings.ToLower(inMsg) == "exit" || strings.ToLower(inMsg) == "quit" {
			fmt.Println("bye!!!")
			os.Exit(0)
		}

		fmt.Println("stream=false")
		callChat(client, inMsg)
		fmt.Println()

		fmt.Println("stream=true")
		streamChat(client, inMsg)
		fmt.Println()
	}

}

func callChat(client deepseek.Client, inMsg string) {
	req := &request.ChatCompletionsRequest{
		Model: "deepseek-chat",
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: inMsg,
			},
		},
		Stream: false,
	}
	resp, err := client.CallChatCompletionsChat(req)
	if err != nil {
		panic(err)
	}
	fmt.Print(resp.Choices[0].Message.Content)
}

func streamChat(client deepseek.Client, inMsg string) {
	req := &request.ChatCompletionsRequest{
		Model: "deepseek-chat",
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: inMsg,
			},
		},
		Stream: true,
	}
	sr, err := client.StreamChatCompletionsChat(req)
	if err != nil {
		panic(err)
	}

	for {
		resp, err := sr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Print(resp.Choices[0].Delta.Content)
	}
}
