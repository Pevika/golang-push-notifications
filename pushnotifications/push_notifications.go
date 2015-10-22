//
// @Author: Geoffrey Bauduin <bauduin.geo@gmail.com>
//

package pushnotifications

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type PushNotification struct {
	sns		*sns.SNS
}

// Create a push notification manager
func NewPushNotification (awsAccessKey string, awsSecretKey string, region string) *PushNotification {
	entity := new(PushNotification)
	cred := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	config := aws.NewConfig().WithRegion(region).WithCredentials(cred)
	entity.sns = sns.New(config)
	return entity
}

// Registers the endpoint into Amazon SNS
func (this *PushNotification) Register (token string, applicationARN string, userData string) (string, error) {
	params := &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(applicationARN),
		Token: aws.String(token),
		Attributes: map[string]*string{
			"Token": aws.String(token),
			"CustomUserData": aws.String(userData),
			"Enabled": aws.String("true"),
		},
		CustomUserData: aws.String(userData),
	}
	resp, err := this.sns.CreatePlatformEndpoint(params)
	if err != nil {
		return "", err
	} else {
		return *resp.EndpointArn, nil
	}
}

// Removes an endpoint from Amazon SNS
func (this *PushNotification) Unregister (arn string) error {
	params := &sns.DeleteEndpointInput{
		EndpointArn: aws.String(arn),
	}
	_, err := this.sns.DeleteEndpoint(params)
	return err
}

// Sends a message to a particular endpoint from Amazon SNS
func (this *PushNotification) Send (arn string, text string) error {
	params := &sns.PublishInput{
		Message: aws.String("{\"default\":\"" + text + "\"}"),
		MessageStructure: aws.String("json"),
		TargetArn: aws.String(arn),
	}
	_, err := this.sns.Publish(params)
	return err
}