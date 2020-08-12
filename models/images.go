package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

// Image stores metadata to be used in posts.
type Image struct {
	PostID   uint
	Filename string
}

// #region Image methods

// Path builds an absolute URL-safe path used to reference this image via web request.
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath builds a path to this image on local disk, relative to the app working directory
func (i *Image) RelativePath() string {
	postID := fmt.Sprintf("%v", i.PostID)
	// ToSlash makes this compatible with windows (stupid backslashes)
	return filepath.ToSlash(filepath.Join("images", "posts", postID, i.Filename))
}

// #endregion

// ImagesService will handle images for the website
type ImagesService interface {
	ByPostID(postID uint) ([]Image, error)
	Create(postID uint, r io.Reader, filename string) error
	Delete(i *Image) error
}

type imagesService struct{}

// NewImagesService is the constructor of ImageService
func NewImagesService() ImagesService {
	return &imagesService{}
}

// ByPostID will get the directory for a post's images, glob it, and return a slice of images.
func (is *imagesService) ByPostID(postID uint) ([]Image, error) {
	path := is.imageDir(postID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Prepare images to return
	images := make([]Image, len(strings))
	for i, imgStr := range strings {
		images[i] = Image{
			Filename: filepath.Base(imgStr),
			PostID:   postID,
		}
	}
	return images, nil
}

// Create will add a new image to a post, storing it locally.
func (is *imagesService) Create(postID uint, r io.Reader, filename string) error {
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

func (is *imagesService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

// #region HELPERS

func (is *imagesService) imageDir(postID uint) string {
	return filepath.Join("public", "images", "posts", fmt.Sprintf("%v", postID))
}

func (is *imagesService) mkImageDir(postID uint) (string, error) {
	postPath := is.imageDir(postID)
	err := os.MkdirAll(postPath, 0755)
	if err != nil {
		return "", err
	}
	return postPath, nil
}

// #endregion
