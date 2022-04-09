package main

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const imageIdLength = 10

type Image struct {
	Id          string
	UserId      string
	Description string
	Location    string
	Size        int64
	CreatedAt   time.Time
}

func NewImage(user *User) *Image {
	return &Image{
		Id:        GenerateId("img", imageIdLength),
		UserId:    user.Id,
		CreatedAt: time.Now(),
	}
}

var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

func (image *Image) CreateFromUrl(imageUrl string) error {
	response, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errImageUrlInvalid
	}

	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	ext, ok := mimeExtensions[mimeType]
	if !ok {
		return errInvalidImageType
	}
	image.Location = image.Id + ext

	dst, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer dst.Close()

	size, err := io.Copy(dst, response.Body)
	if err != nil {
		return err
	}
	image.Size = size

	return globalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return err
	}

	mimeType := http.DetectContentType(buf)
	ext, ok := mimeExtensions[mimeType]
	if !ok {
		return errInvalidImageType
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	image.Location = image.Id + ext

	dst, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return err
	}
	image.Size = size

	return globalImageStore.Save(image)
}
