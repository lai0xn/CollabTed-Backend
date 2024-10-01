package cloudinary

import (
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloud cloudinary.Cloudinary

func Connect() {
	cloud, err := cloudinary.New()
	if err != nil {
		logger.Logger.Err(err)
	}
	cloud.Config.URL.Secure = true
	logger.Logger.Info().Msg("Connected to cloudinary")
}

func GetUploader() *uploader.API {
	return &cloud.Upload
}
