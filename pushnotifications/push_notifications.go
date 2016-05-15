//
// @Author: Geoffrey Bauduin <bauduin.geo@gmail.com>
//

package pushnotifications

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
    "encoding/json"
)

type PushNotification struct {
	sns		*sns.SNS
}

type Alert struct {
	Body			*string			`json:"body"`
	LocKey			*string			`json:"loc-key"`
	LocArgs			*[]interface{}	`json:"loc-args"`
	ActionLocKey	*string			`json:"action-loc-key"`
}

type Push struct {
	Alert       *string         `json:"alert,omitempty"`
    Sound       *string         `json:"sound,omitempty"`
    Data        interface{}     `json:"custom_data"`
	Badge		*int			`json:"badge,omitempty"`
}

type wrapper struct {
    APNS        string         	`json:"APNS"`
    APNSSandbox string         	`json:"APNS_SANDBOX"`
    Default     string         	`json:"default"`
	GCM			string			`json:"GCM"`      
}

type iosPush struct {
    APS         Push           	`json:"aps"`
}

type gcmPush struct {
	Message		string			`json:"message"`
	Custom		interface{}		`json:"custom"`
}

type gcmPushWrapper struct {
	Data		gcmPush			`json:"data"`	
}

// Create a push notification manager
func NewPushNotification (awsAccessKey string, awsSecretKey string, region string) *PushNotification {
	entity := new(PushNotification)
	cred := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	config := aws.NewConfig().WithRegion(region).WithCredentials(cred)
	sess := session.New(config)
	entity.sns = sns.New(sess)
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
func (this *PushNotification) Send (arn string, data *Push) error {
    msg := wrapper{}
    ios := iosPush{
        APS: *data,
    }
    b, err := json.Marshal(ios)
    if err != nil {
        return err
    }
    msg.APNS = string(b[:])
    msg.APNSSandbox = msg.APNS
    msg.Default = *data.Alert
	gcm := gcmPushWrapper{
		Data: gcmPush{
			Message: *data.Alert,
			Custom: data.Data,
		},
	}
	b, err = json.Marshal(gcm)
	if err != nil {
		return err
	}
	msg.GCM = string(b[:])
    pushData, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    m := string(pushData[:])
	params := &sns.PublishInput{
		Message: aws.String(m),
		MessageStructure: aws.String("json"),
		TargetArn: aws.String(arn),
	}
	_, err = this.sns.Publish(params)
	return err
}