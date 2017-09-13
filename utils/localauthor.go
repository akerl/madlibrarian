package utils

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type localAuthor struct {
}

func (la localAuthor) Upload(s *Story, bucket, prefix string) (*Story, error) {
	chunks, err := la.chunks(s)
	if err != nil {
		return &Story{}, err
	}
	var chunkNames []string
	for chunk, list := range chunks {
		chunkNames = append(chunkNames, chunk)
		if err := la.uploadChunk(bucket, prefix, chunk, list); err != nil {
			return &Story{}, err
		}
	}

	newStory := Story{
		Meta: Metadata{
			Type:     "s3",
			Template: s.Meta.Template,
		},
		Data: map[string]interface{}{
			"s3": map[string]string{
				"bucket": bucket,
				"prefix": prefix,
			},
			"chunks": chunkNames,
		},
	}
	return &newStory, nil
}

func (la localAuthor) uploadChunk(bucket, prefix, chunk string, list []string) error {
	counter := -1
	for _, line := range list {
		counter++
		if err := la.uploadItem(bucket, prefix, chunk, counter, line); err != nil {
			return err
		}
	}
	key := makeS3Key(prefix, chunk, "max")
	counterStr := fmt.Sprintf("%d", counter)
	return la.uploadFile(bucket, key, counterStr)
}

func (la localAuthor) uploadItem(bucket, prefix, chunk string, counter int, line string) error {
	key := makeS3Key(prefix, chunk, counter)
	return la.uploadFile(bucket, key, line)
}

func (la localAuthor) uploadFile(bucket, key, body string) error {
	uploader := s3Session.Uploader()
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(body),
	})
	return err
}

func (la localAuthor) Funcs(s *Story) (template.FuncMap, error) {
	funcMap := template.FuncMap{}
	chunks, err := la.chunks(s)
	if err != nil {
		return funcMap, err
	}
	for chunk, list := range chunks {
		funcMap[chunk] = localRandomizer(list)
	}
	return funcMap, nil
}

func (la localAuthor) chunks(s *Story) (map[string][]string, error) {
	result := make(map[string][]string)
	for name, iface := range s.Data {
		list, err := ifaceToStringSlice(iface)
		if err != nil {
			return result, err
		}
		result[name] = list
	}
	return result, nil
}

func localRandomizer(list []string) func() (string, error) {
	return func() (string, error) {
		index, err := getRand(len(list) - 1)
		if err != nil {
			return "", err
		}
		return list[index], nil
	}
}
