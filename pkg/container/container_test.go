package container

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/prometheus/common/log"
	"os"
	"testing"
)

func TestCopyImagesToBackUpRegistry(t *testing.T) {
	image, err := crane.Pull("containous/whoami:latest")
	if err != nil{
		log.Error(err)
		os.Exit(1)
	}
    err = crane.Push(image, "quay.io/cargaona/whoami", crane.WithAuthFromKeychain(authn.DefaultKeychain))
    if err != nil {
    	log.Error(err)
    	os.Exit(1)
	}
	os.Exit(2)
//	err := crane.Copy("nginx:latest","quay.io/cargaona/nginx", crane.WithAuthFromKeychain(authn.DefaultKeychain))
//	if err != nil {
//		log.Error(err)
//		os.Exit(1)
//	}

	repo, err := name.NewRepository("quay.io/cargaona/whoamibackup")
	if err != nil {
		fmt.Println(repo)
		log.Error(err)
		os.Exit(1)
	}

	// Get original image

	oldImageName, err := name.ParseReference("nginx")
	if err != nil {
		log.Info("Salió acá")
	}
	oldRemoteImage, err := remote.Image(oldImageName)
	//repo, _:= name.NewRepository(fmt.Sprintf("quay.io/cargaona/%s", newImageName))
	if err != nil {
		os.Exit(1)
	}
	oldRemoteImage.Manifest()
	// create new image
	newImageName, err := name.ParseReference("quay.io/cargaona/whoamibackup")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	newRemoteImage, err := remote.Image(newImageName, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	tag, _ := name.NewTag("Latest")
	err = remote.Tag(tag, newRemoteImage)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = remote.Write(newImageName, newRemoteImage)

	if err != nil {
		log.Error(err)
	}
}
