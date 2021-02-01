package container

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "k8s.io/api/core/v1"
	"strings"
)

func CheckImagesSource(ctx context.Context, images v1.PodSpec, backupRegistry string) (bool, map[int]string) {
	imagesToBackup := make(map[int]string)
	for index, container := range images.Containers {
		if imageFromBackupRegistry(container.Image, backupRegistry) == false {
			imagesToBackup[index] = container.Image
		}
	}
	// Nothing to backup
	if len(imagesToBackup) == 0 {
		return false, imagesToBackup
	}

	return true, imagesToBackup
}

func imageFromBackupRegistry(image string, backupRegistry string) bool {
	//TODO: improve this validation.
	if strings.Contains(image, backupRegistry) {
		return true
	}
	return false
}

func CopyImagesToBackUpRegistry(ctx context.Context, images map[int]string, backupRegistry string) (map[int]string, error) {
	//TODO: Support for schema v1 images and put more attention on the tags.

	newImages := make(map[int]string)
	for key, imageName := range images {
		image, err := crane.Pull(imageName)
		if err != nil {
			return nil, err
		}
		sanitizedName := getCleanImageName(imageName)
		err = crane.Push(image, fmt.Sprintf("%s/%s", backupRegistry, sanitizedName), crane.WithAuthFromKeychain(authn.DefaultKeychain))
		if err != nil {
			return nil, err
		}
		newImages[key] = fmt.Sprintf("%s/%s", backupRegistry, sanitizedName)
	}
	return newImages, nil
}

func getCleanImageName(imageName string) string {
	sanitizedName := strings.Split(imageName, "/")
	return sanitizedName[len(sanitizedName)-1]
}

