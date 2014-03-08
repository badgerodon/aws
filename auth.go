package aws

import (
	"errors"
	"os"
)

type Auth struct {
	AccessKeyID, SecretAccessKey string
}

func EnvAuth() (auth Auth, err error) {
	auth.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	if auth.AccessKeyID == "" {
		auth.AccessKeyID = os.Getenv("AWS_ACCESS_KEY")
	}

	auth.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if auth.SecretAccessKey == "" {
		auth.SecretAccessKey = os.Getenv("AWS_SECRET_KEY")
	}
	if auth.AccessKeyID == "" {
		err = errors.New("AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment")
	}
	if auth.SecretAccessKey == "" {
		err = errors.New("AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY not found in environment")
	}
	return
}
