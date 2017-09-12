package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	keyMax = 100000000
)

var s3Session = s3SessionObj{}

type s3SessionObj struct {
	session    *session.Session
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	client     *s3.S3
}

func (s *s3SessionObj) loadSession() {
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	s.session = session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func (s *s3SessionObj) Downloader() *s3manager.Downloader {
	if s.session == nil {
		s.loadSession()
	}
	if s.downloader == nil {
		s.downloader = s3manager.NewDownloader(s.session)
	}
	return s.downloader
}

func (s *s3SessionObj) Uploader() *s3manager.Uploader {
	if s.session == nil {
		s.loadSession()
	}
	if s.uploader == nil {
		s.uploader = s3manager.NewUploader(s.session)
	}
	return s.uploader
}

func (s *s3SessionObj) Client() *s3.S3 {
	if s.session == nil {
		s.loadSession()
	}
	if s.client == nil {
		s.client = s3.New(s.session)
	}
	return s.client
}

func makeS3Key(prefix, chunk string, index int) string {
	return fmt.Sprintf("%s/%s/%09d", prefix, chunk, index)
}

// Upload takes a Story and puts its quotes in S3
func Upload(s Story, bucket string, prefix string) error {
	uploader := s3Session.Uploader()

	for chunk, iface := range s.Data {
		list, err := ifaceToStringSlice(iface)
		if err != nil {
			return err
		}
		step := keyMax / len(list)
		counter := 0
		for _, line := range list {
			counter += step
			key := makeS3Key(prefix, chunk, counter)
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
