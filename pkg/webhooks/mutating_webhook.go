package webhooks

import (
	"context"
	"encoding/json"
	"github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/container"
	corev1 "k8s.io/api/core/v1"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodImageMutator struct {
	Client  client.Client
	decoder *admission.Decoder
	Config  configuration.Config
}

func (a *PodImageMutator) Handle(ctx context.Context, req admission.Request) admission.Response {

	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	imageToBackupExist, images := container.CheckImagesSourceFromPod(pod, a.Config.BackupRegistry)
	if imageToBackupExist == false {
		pod.Annotations["image-cloner-controller.backed-images"] = "true"
		return admission.Response{}
	}

	newImages, err := container.CopyImagesToBackUpRegistry(ctx, images, a.Config.BackupRegistry)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	// Change image values on pod.
	for key, value := range newImages {
		pod.Spec.Containers[key].Image = value
	}

	pod.Annotations["image-cloner-controller.backed-images"] = "true"

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (a *PodImageMutator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
