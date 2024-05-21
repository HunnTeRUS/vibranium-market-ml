package order_queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sync/atomic"
)

func (d *Disruptor) SaveSnapshotToS3() error {
	bucket := os.Getenv("S3_BUCKET")
	keyPrefix := os.Getenv("S3_KEY_PREFIX")
	if bucket == "" || keyPrefix == "" {
		return errors.New("S3_BUCKET or S3_KEY_PREFIX environment variable is not set")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(d.buffer)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	svc := s3.New(sess)
	key := fmt.Sprintf("%s/snapshot-%d.gob", keyPrefix, atomic.LoadInt64(&d.writeCursor))
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	return err
}

func (d *Disruptor) LoadSnapshotFromS3() error {
	bucket := os.Getenv("S3_BUCKET")
	keyPrefix := os.Getenv("S3_KEY_PREFIX")
	if bucket == "" || keyPrefix == "" {
		return errors.New("S3_BUCKET or S3_KEY_PREFIX environment variable is not set")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	svc := s3.New(sess)

	listObjectsOutput, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keyPrefix),
	})
	if err != nil {
		return err
	}

	if len(listObjectsOutput.Contents) == 0 {
		return errors.New("no snapshots found in S3")
	}

	for _, object := range listObjectsOutput.Contents {
		_, err := svc.PutObjectTagging(&s3.PutObjectTaggingInput{
			Bucket: aws.String(bucket),
			Key:    object.Key,
			Tagging: &s3.Tagging{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("in-use"),
						Value: aws.String("true"),
					},
				},
			},
		})
		if err != nil {
			continue
		}

		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    object.Key,
		})
		if err != nil {
			return err
		}
		defer result.Body.Close()

		decoder := gob.NewDecoder(result.Body)
		err = decoder.Decode(&d.buffer)
		if err != nil {
			return err
		}

		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    object.Key,
		})
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("no available snapshots found in S3")
}
