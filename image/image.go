package image

import "net/textproto"

const (
	myBucket        = "bucketpashakrut"
	accessKeyID     = "AKIAX43I6M7I7SFYYAWM"
	secretAccessKey = "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"
	region          = "eu-central-1"
)

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}
