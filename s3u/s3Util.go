package s3u

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"time"
)

var (
	timeout time.Duration = 10 * time.Second
)

type S3Util struct {
	region          string
	accessKey       string
	secretAccessKey string
	svc             *s3.S3
}

func NewS3Util(region string, accessKey string, secretAccessKey string) *S3Util {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession())
	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	config := &aws.Config{
		Region: &region,
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretAccessKey,
		}),
	}
	config.WithCredentialsChainVerboseErrors(true)
	svc := s3.New(sess, config)
	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	defer cancelFn()
	return &S3Util{region: region, accessKey: accessKey, secretAccessKey: secretAccessKey, svc: svc}
}

func (s3u *S3Util) Presign(key string, bucket string) string {
	sdkReq, _ := s3u.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	var u string
	var err1 error
	u, err1 = sdkReq.Presign(1 * time.Minute)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println("url:", u)
	return u
}

func (s3u *S3Util) List(bucket string) {
	listObjectsInput := new(s3.ListObjectsInput)
	listObjectsInput = listObjectsInput.SetBucket(bucket)
	otList, err := s3u.svc.ListObjects(listObjectsInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, o := range otList.Contents {
		_ = i
		/*if i > 0 {
			break
		}*/
		fmt.Println(aws.StringValue(o.Key))
		/*key := aws.StringValue(o.Key)
		sdkReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		var u string
		var err1 error
		var signedHeaders http.Header
		//u, signedHeaders, err1 = sdkReq.PresignRequest(1 * time.Minute)
		u, err1 = sdkReq.Presign(1 * time.Minute)
		if err1 != nil {
			fmt.Println(err1)
		}
		fmt.Println("key:", u, "signedHeaders:", signedHeaders)*/

	}

	//req,outPut := svc.ListObjectsRequest(listObjectsInput)

}

//上传文件
func (s3u *S3Util) uploadFile(file string, bucket string, ctx context.Context) {
	f, _ := os.Open("D:/github/golang/src/AwsSdk/main.go")
	defer f.Close()
	_, err := s3u.svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("main.go"),
		Body:   S3Body{f},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("successfully uploaded file to ", bucket)
}

type S3Body struct {
	file *os.File
}

func (b S3Body) Read(p []byte) (n int, err error) {
	return b.file.Read(p)
}

func (b S3Body) Seek(offset int64, whence int) (int64, error) {
	return b.file.Seek(offset, whence)
}
