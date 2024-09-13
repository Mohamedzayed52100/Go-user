package repository

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) SendMessage(ctx context.Context, req *userProto.SendMessageRequest) (*userProto.SendMessageResponse, error) {
	if len(req.GetTo()) == 0 {
		if err := r.SendBroadcastMessage("broadcast", req.GetMessage()); err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	} else {
		for _, v := range req.GetTo() {
			var email string
			if err := r.SharedDbConnection.Model(&user.User{}).Where("id = ?", v).Select("email").Scan(&email).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			if err := r.SendPersonalMessage(
				ctx,
				email,
				req.GetMessage(),
			); err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}
	}

	return &userProto.SendMessageResponse{
		Code:    http.StatusOK,
		Message: "Message sent successfully",
	}, nil
}

func createQueue(sess *session.Session, queueName string) (*sqs.CreateQueueOutput, error) {
	sqsClient := sqs.New(sess)
	result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]*string{
			"DelaySeconds":      aws.String("0"),
			"VisibilityTimeout": aws.String("60"),
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func getQueueURL(sess *session.Session, queue string) (string, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}

func (r *UserRepository) SendPersonalMessage(ctx context.Context, queueName, messageBody string) error {
	var queueUrl string

	queueName = strings.ReplaceAll(queueName, ".", "-")
	queueName = strings.ReplaceAll(queueName, "@", "-")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return err
	}

	sqsClient := sqs.New(sess)
	if queueUrl, err = getQueueURL(sess, queueName); err != nil {
		if _, err := createQueue(sess, queueName); err != nil {
			return err
		}

		queueUrl, err = getQueueURL(sess, queueName)
		if err != nil {
			return err
		}
	}

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	id := uuid.New().String()
	fullMessageBody := fmt.Sprintf(`{"id": "%s", "from": {"firstName": "%s",  "lastName": "%s", "avatar": "%s"}, "body": "%s", "type": "message", "seen": false, "createdAt": "%v"}`, id, currentUser.FirstName, currentUser.LastName, currentUser.Avatar, messageBody, time.Now().UTC())

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(fullMessageBody),
	})

	return err
}

func (r *UserRepository) SendBroadcastMessage(queueName, messageBody string) error {
	var queueUrl string

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return err
	}

	sqsClient := sqs.New(sess)

	if queueUrl, err = getQueueURL(sess, queueName); err != nil {
		if _, err := createQueue(sess, queueName); err != nil {
			return err
		}

		queueUrl, err = getQueueURL(sess, queueName)
		if err != nil {
			return err
		}
	}

	id := uuid.New().String()
	fullMessageBody := fmt.Sprintf(`{"id": "%s", "from": "Support", "body": "%s", "type": "broadcast", "seen": false, "createdAt": "%v"}`, id, messageBody, time.Now().UTC())

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(fullMessageBody),
	})

	return err
}
