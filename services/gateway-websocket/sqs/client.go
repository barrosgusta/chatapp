package sqs

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Producer is an interface for sending messages to SQS.
type Producer interface {
	SendMessage(ctx context.Context, body string) error
}

type SQSProducer struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSProducer(ctx context.Context) *SQSProducer {
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
	// cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	
	// For production, use the environment variables
	cfg, err := config.LoadDefaultConfig(ctx,
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

	client := sqs.NewFromConfig(cfg)
	return &SQSProducer{client: client, queueURL: queueURL}
}

func (p *SQSProducer) SendMessage(ctx context.Context, body string) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(body),
	})
	return err
}
