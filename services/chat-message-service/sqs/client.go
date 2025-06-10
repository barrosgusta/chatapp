package sqs

import (
	"chat-message-service/dynamodb"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Consumer struct {
    client  *sqs.Client
    queueURL string
}

func NewSQSConsumer() *Consumer {
    queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL is not set")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		log.Fatal("AWS_ACCESS_KEY_ID is not set")
	}
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		log.Fatal("AWS_SECRET_ACCESS_KEY is not set")
	}

	// For local development
	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	
	// For production, use the environment variables
	cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(
        aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
            accessKey, secretKey, "",
        )),
    ),
    config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
    return &Consumer{
        client:  sqs.NewFromConfig(cfg),
        queueURL: queueURL,
    }
}

func (c *Consumer) StartConsuming() {
    for {
        output, err := c.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(c.queueURL),
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     10,
        })
        if err != nil {
            log.Printf("SQS receive error: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }

        for _, msg := range output.Messages {
            var chatMsg dynamodb.ChatMessage
            if err := json.Unmarshal([]byte(*msg.Body), &chatMsg); err != nil {
                log.Printf("Failed to unmarshal SQS message: %v", err)
                continue
            }
            if err := dynamodb.SaveMessage(chatMsg); err != nil {
                log.Printf("Failed to save message to DynamoDB: %v", err)
                continue
            } else {
                log.Printf("Message saved to DynamoDB: %s", chatMsg.Id)
            }
            _, err = c.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      aws.String(c.queueURL),
                ReceiptHandle: msg.ReceiptHandle,
            })
            if err != nil {
                log.Printf("Failed to delete SQS message: %v", err)
            }
        }
    }
}