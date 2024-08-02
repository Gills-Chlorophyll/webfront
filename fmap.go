package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

/* This will house all the function maps
funcion maps are the function references that get passed into templates for advanced logic on the template
*/

// isNotEmptyString: to be used in templates to see if the value is not empty string
func isNotEmptyString(a string) bool {
	return a != ""
}

func countToRange(count int) []int {
	result := []int{}
	for i := 1; i <= count; i++ {
		result = append(result, i)
	}
	return result
}
func presignImageUrl(key string) string {
	req, _ := AWS_S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(VULTR_S3_BUCKET),
		Key:    aws.String(key),
	})

	// Set the URL to expire in one hour
	urlStr, err := req.Presign(20 * time.Minute)
	if err != nil {
		fmt.Println("Error presigning request:", err)
		return ""
	}

	return urlStr
}
