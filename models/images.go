package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImageService will handle images for the website
type ImageService interface {
	Create(postID uint, r io.Reader, filename string) error
}

type imageService struct{}

// NewImageService is the constructor of ImageService
func NewImageService() ImageService {
	return &imageService{}
}

// Create will add a new image to a post, storing it locally.
func (is *imageService) Create(postID uint, r io.Reader, filename string) error {
	path, err := is.mkImagePath(postID)
	if err != nil {
		return err
	}
	// Create destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) mkImagePath(postID uint) (string, error) {
	postPath := filepath.Join("images", "posts", fmt.Sprintf("%v", postID))
	err := os.MkdirAll(postPath, 0755)
	if err != nil {
		return "", err
	}
	return postPath, nil
}