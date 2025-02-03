package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
)

func main() {
	// create client
	client, _ := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	// get user input message
	fmt.Println("Demo of DeepSeek Go client calling chat api with model=deepseek-reasoner, stream=true")
	fmt.Printf("\nMessage >>> ")
	reader := bufio.NewReader(os.Stdin)
	lineBytes, _, _ := reader.ReadLine()
	inMsg := string(lineBytes)

	// send request to go deepseek client
	req := &request.ChatCompletionsRequest{
		Model: deepseek.DEEPSEEK_REASONER_MODEL,
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: inMsg,
			},
		},
		Stream: true,
	}
	sr, _ := client.StreamChatCompletionsReasoner(context.Background(), req)

	// process response
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
	fmt.Println()
}
