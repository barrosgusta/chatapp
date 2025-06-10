package dynamodb

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *dynamodb.Client

type ChatMessage struct {
	Id        string `json:"id"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func Init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

func SaveMessage(msg ChatMessage) error {
	av, err := attributevalue.MarshalMap(msg)
	if err != nil {
		return err
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("CHAT_MESSAGES_TABLE")),
		// For local development, you can use a hardcoded table name
		// TableName: aws.String("chat_messages"),
		Item:      av,
	})
	return err
}

func GetRecentMessages() ([]ChatMessage, error) {
	out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("CHAT_MESSAGES_TABLE")),
		// For local development, you can use a hardcoded table name
		// TableName: aws.String("chat_messages"),
		Limit:     aws.Int32(50),
	})
	if err != nil {
		return nil, err
	}
	var messages []ChatMessage
	err = attributevalue.UnmarshalListOfMaps(out.Items, &messages)
	return messages, err
}