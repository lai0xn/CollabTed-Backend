package cloudinary

import (
	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloud *cloudinary.Cloudinary

func Connect() {
	var err error
	cloud, err = cloudinary.NewFromURL(config.CLOUDINARY_URL)
	if err != nil {
		logger.Logger.Err(err).Msg("Failed to connect to Cloudinary")
		return
	}

	cloud.Config.URL.Secure = true
	logger.Logger.Info().Msg("Connected to Cloudinary")
}

func GetUploader() *uploader.API {
	return &cloud.Upload
}
