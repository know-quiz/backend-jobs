package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	"github.com/google/uuid"
)

// ChatGPTResponse simulates the JSON response from ChatGPT
type ChatGPTResponse struct {
	Questions []struct {
		Question string `json:"question"`
		Option1  string `json:"option1"`
		Option2  string `json:"option2"`
		Option3  string `json:"option3"`
		Option4  string `json:"option4"`
		Answer   string `json:"answer"`
	} `json:"questions"`
}

// fetchChatGPTData generates data from ChatGPT
func fetchChatGPTData() (*ChatGPTResponse, error) {
	// Simulate API request to ChatGPT
	// In a real scenario, you would use fasthttp to make a POST request to the ChatGPT API endpoint
	// and pass the necessary headers and body (API key, prompt, etc.)
	//
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "You are a random quiz generator for kids. Now generate a an array of 20 questions and four options for each question in json format. Make sure you don't return anything other than the json response. Also add a field to indiciate the right answer which can be used to render a UI and show the quiz to the kid. The response structure should have the fields array of question, option1, option2, option3, option4, answer. all in lower case. The answer field should have the name of the option that has the correct answer. example if option1 is correct answer for a question the answer field should be 'answer':'option1' ",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}
	response := &ChatGPTResponse{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), response)
	if err != nil {
		return nil, err
	}
	for i, question := range response.Questions {
		fmt.Println(i, question.Question)
	}
	return response, nil

}

func writeToFireStore(response *ChatGPTResponse) {
	ctx := context.Background()
	d, _ := base64.StdEncoding.DecodeString(os.Getenv("GCP_CREDS_JSON_BASE64"))
	//sa := option.WithCredentialsFile("/Users/dileep/Downloads/firestore-quiz-app-415922-d27bce37d853.json")
	client, err := firestore.NewClient(ctx, "quiz-app-415922", option.WithCredentialsJSON(d))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Specify the collection and document
	collection := "quiz-questions"

	// Write data to Firestore
	for _, question := range response.Questions {
		_, err = client.Collection(collection).Doc(uuid.NewString()).Set(ctx, question)
		if err != nil {
			log.Fatalf("Failed to write data to Firestore: %v", err)
		}
	}
	fmt.Println("Data successfully written to Firestore")
}

func main() {
	response, err := fetchChatGPTData()
	if err != nil {
		log.Fatalf("Failed to fetch data from ChatGPT: %v", err)
	}
	writeToFireStore(response)
}
