package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	keyMax = 100000000
)

// Upload takes a Story and puts its quotes in S3
func Upload(s Story, bucket string, prefix string) error {
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
	uploader := s3manager.NewUploader(sess)

	for chunk, iface := range s.Data {
		list, err := ifaceToStringSlice(iface)
		if err != nil {
			return err
		}
		// TODO: Check this math
		step := keyMax / len(list)
		counter := 0
		for _, line := range list {
			counter += step
			key := fmt.Sprintf("%s/%s/%09d", prefix, chunk, counter)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   strings.NewReader(line),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
