package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImageService will handle images for the website
type ImageService interface {
	ByPostID(postID uint) ([]string, error)
	Create(postID uint, r io.Reader, filename string) error
}

type imageService struct{}

// NewImageService is the constructor of ImageService
func NewImageService() ImageService {
	return &imageService{}
}

// ByPostID will get the directory for a post's images and glob it.
func (is *imageService) ByPostID(postID uint) ([]string, error) {
	path := is.imageDir(postID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	return strings, nil
}

// Create will add a new image to a post, storing it locally.
func (is *imageService) Create(postID uint, r io.Reader, filename string) error {
	path, err := is.mkImageDir(postID)
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

// #region HELPERS

func (is *imageService) imageDir(postID uint) string {
	return filepath.Join("images", "posts", fmt.Sprintf("%v", postID))
}

func (is *imageService) mkImageDir(postID uint) (string, error) {
	postPath := is.imageDir(postID)
	err := os.MkdirAll(postPath, 0755)
	if err != nil {
		return "", err
	}
	return postPath, nil
}

// #endregion