package webhooks

import (
	"context"
	"fmt"
	"github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/container"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodImageValidator struct {
	Client  client.Client
	decoder *admission.Decoder
	Config  configuration.Config
}

func (a *PodImageValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	imagesToBackupExist, imagesToBackup := container.CheckImagesSourceFromPod(pod, a.Config.BackupRegistry)

	if pod.Annotations["image-cloner-controller.backed-images"] != "true" {
		return admission.Denied(fmt.Sprintf("Missing annotation for pod/%s", pod.Name))
	}

	if imagesToBackupExist != false {
		return admission.Denied(fmt.Sprintf("There are images to backup %v", imagesToBackup))
	}

	return admission.Allowed("")
}

func (a *PodImageValidator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}