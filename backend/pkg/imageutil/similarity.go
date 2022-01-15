package imageutil

import (
	"context"
	"image"

	"github.com/vitali-fedulov/images/v2"
	"golang.org/x/sync/errgroup"
)

func IsSimilarImage(imageA image.Image, imageB image.Image) bool {
	hashA, imgSizeA := images.Hash(imageA)
	hashB, imgSizeB := images.Hash(imageB)

	return images.Similar(hashA, hashB, imgSizeA, imgSizeB)
}

func IsSimilarImageByURLs(ctx context.Context, imageAURL string, imageBURL string) (bool, error) {
	eg := errgroup.Group{}
	var imageA, imageB image.Image
	eg.Go(func() error {
		var err error
		imageA, err = DownloadImage(ctx, imageAURL)
		return err
	})
	eg.Go(func() error {
		var err error
		imageB, err = DownloadImage(ctx, imageBURL)
		return err
	})
	if err := eg.Wait(); err != nil {
		return false, err
	}

	return IsSimilarImage(imageA, imageB), nil
}
