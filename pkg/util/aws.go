package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() *session.Session {
	return NewSessionWithRegion("us-west-2")
}

func NewSessionWithRegion(region string) *session.Session {
	config := aws.NewConfig().WithRegion(region)
	return session.Must(session.NewSession(config))
}
