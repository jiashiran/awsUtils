package main

import (
	"fmt"
	"awsUtils/s3u"
)

var (
	region          	= ""
	accessKey              	= ""
	secretAccessKey        	= ""
	bucket  	        = ""
)

func main() {

	s3 := s3u.NewS3Util(region,accessKey,secretAccessKey)
	//list(svc)
	s3.Presign("key","bucket")
	fmt.Println("over")
}
