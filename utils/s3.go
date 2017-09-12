package utils

import (
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Story struct {
}

func (ss s3Story) Funcs(s *Story) (template.FuncMap, error) {
	funcMap := template.FuncMap{}

	s3InfoIface, ok := s.Data["s3"]
	if !ok {
		return funcMap, fmt.Errorf("No S3 section in YAML")
	}
	s3InfoRaw, ok := s3InfoIface.(map[interface{}]interface{})
	if !ok {
		return funcMap, fmt.Errorf("S3 section is not a map of strings to strings")
	}

	s3Info := make(map[string]string)
	for key, value := range s3InfoRaw {
		sKey, ok := key.(string)
		if !ok {
			return funcMap, fmt.Errorf("S3 key not a string: %s", key)
		}
		sValue, ok := value.(string)
		if !ok {
			return funcMap, fmt.Errorf("S3 value not a string: key %s", key)
		}
		s3Info[sKey] = sValue
	}

	chunksIface, ok := s.Data["chunks"]
	if !ok {
		return funcMap, fmt.Errorf("Chunks not defined")
	}
	chunks, err := ifaceToStringSlice(chunksIface)
	if err != nil {
		return funcMap, err
	}

	for _, name := range chunks {
		funcMap[name] = s3Randomizer(s3Info, name)
	}
	return funcMap, nil
}

func s3Randomizer(s3Info map[string]string, chunk string) func() (string, error) {
	return func() (string, error) {
		bucket, ok := s3Info["bucket"]
		if !ok {
			return "", fmt.Errorf("No bucket in config")
		}
		prefix, ok := s3Info["prefix"]
		if !ok {
			return "", fmt.Errorf("No prefix in config")
		}

		client := s3Session.Client()
		random, err := getRandomIndex(keyMax)
		if err != nil {
			return "", err
		}
		randKey := makeS3Key(prefix, chunk, random)

		list, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(bucket),
			MaxKeys:    aws.Int64(1),
			StartAfter: aws.String(randKey),
		})
		if err != nil {
			return "", nil
		}
		if len(list.Contents) != 1 {
			return "", fmt.Errorf("No results found in bucket for %s", chunk)
		}
		key := list.Contents[0].Key

		resultObj, err := client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    key,
		})
		if err != nil {
			return "", err
		}

		result, err := ioutil.ReadAll(resultObj.Body)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
}
