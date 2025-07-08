package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(filePath string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	resp, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	fmt.Println("Upload successful, URL:", resp.SecureURL)
	return resp.SecureURL, nil
}
