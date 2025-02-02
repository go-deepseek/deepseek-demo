package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/config"
	"github.com/go-deepseek/deepseek/request"
	"github.com/go-deepseek/deepseek/response"
)

const DEEPSEEK_API_KEY = `sk-123cd456c78d9be0b123de45cf6789b0` // replace with valid one

func main() {
	config := config.Config{
		ApiKey:         DEEPSEEK_API_KEY,
		TimeoutSeconds: 120,
	}
	client, err := deepseek.NewClientWithConfig(config)
	if err != nil {
		panic(err)
	}

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
		callChat(client, "deepseek-chat", inMsg)
		fmt.Println()

		fmt.Println("stream=true")
		streamChat(client, "deepseek-chat", inMsg)
		fmt.Println()

		// fmt.Println("stream=false")
		// callChat(client, "deepseek-reasoner", inMsg)
		// fmt.Println()

		// fmt.Println("stream=true")
		// streamChat(client, "deepseek-reasoner", inMsg)
		// fmt.Println()
	}

}

func callChat(client deepseek.Client, model, inMsg string) {
	req := &request.ChatCompletionsRequest{
		Model: model,
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: inMsg,
			},
		},
		Stream: false,
	}

	var resp *response.ChatCompletionsResponse
	var err error
	if model == "deepseek-chat" {
		ctx := context.Background()
		// ctx, _ = context.WithTimeout(ctx, time.Second*2)
		resp, err = client.CallChatCompletionsChat(ctx, req)
	} else {
		resp, err = client.CallChatCompletionsReasoner(context.Background(), req)
	}
	if err != nil {
		panic(err)
	}

	fmt.Print(resp.Choices[0].Message.Content)
	if model == "deepseek-reasoner" {
		fmt.Println()
		fmt.Print(resp.Choices[0].Message.Content)
	}
}

func streamChat(client deepseek.Client, model, inMsg string) {
	req := &request.ChatCompletionsRequest{
		Model: model,
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: inMsg,
			},
		},
		Stream: true,
	}

	var sr response.StreamReader
	var err error
	if model == "deepseek-chat" {
		sr, err = client.StreamChatCompletionsChat(context.Background(), req)
	} else {
		sr, err = client.StreamChatCompletionsReasoner(context.Background(), req)
	}
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
		if resp.Choices[0].Delta.Content != "" {
			fmt.Print(resp.Choices[0].Delta.Content)
		} else {
			fmt.Print(resp.Choices[0].Delta.ReasoningContent)
		}
	}
}
