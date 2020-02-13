package image

import (
	"mime/multipart"

	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

func HandleUpload(repo users.UserRepo, file multipart.File, header *multipart.FileHeader, userID int) (URL string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return "", errors.Wrap(err, "error creating SDK session")
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(header.Filename),
		Body:   file,
	})
	if err != nil {
		return "", errors.Wrap(err, "error upload file to AWS s3")
	}
	err = repo.UpdateAvatarURL(result.Location, userID)
	if err != nil {
		return "", errors.Wrap(err, "error update field URL of users in DB")
	}
	return result.Location, nil
}
