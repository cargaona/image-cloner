package container

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "k8s.io/api/core/v1"
	"strings"
)

const backupRegistry = "quay.io/cargaona"

func CheckImagesSource(ctx context.Context, images v1.PodSpec) (bool, map[int]string) {
	imagesToBackup := make(map[int]string)
	for index, container := range images.Containers {
		if imageFromBackupRegistry(container.Image) == false {
			imagesToBackup[index] = container.Image
		}
	}
	// Nothing to backup
	if len(imagesToBackup) == 0 {
		return false, imagesToBackup
	}

	return true, imagesToBackup
}

func imageFromBackupRegistry(image string) bool {
	//Can be a better validation
	if strings.Contains(image, backupRegistry) {
		return true
	}
	return false
}

func CopyImagesToBackUpRegistry(ctx context.Context, images map[int]string) (map[int]string, error) {
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

func ValidateRedeployedDaemonset(ctx context.Context, status int32, deployedImages v1.PodSpec, mustHaveImages map[int]string) error {
	if status > 0 {
		return fmt.Errorf("there are unavailable daemonsets")
	}
	for key, value := range mustHaveImages {
		if deployedImages.Containers[key].Image == value {
			continue
		} else {
			return fmt.Errorf("image %s from container %s was not redeployed succesfully", value, deployedImages.Containers[key].Name)
		}
	}
	return nil
}
