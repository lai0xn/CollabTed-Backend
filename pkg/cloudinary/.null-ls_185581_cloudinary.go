package cloudinary

import "github.com/cloudinary/cloudinary-go/v2"

var cloud = cloudinary.Cloudinary

func New() {
	cid, err := cloudinary.New()
	if err != nil {
		panic(err)
	}
}
