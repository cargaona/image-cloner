package webhooks

import (
	"context"
	"fmt"
	"github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/container"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type DeploymentImageValidator struct {
	Client  client.Client
	decoder *admission.Decoder
	Config  configuration.Config
	Logger  logr.Logger
}

func (a *DeploymentImageValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	instance := &appsv1.Deployment{}

	err := a.decoder.Decode(req, instance)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	imagesToBackupExist, imagesToBackup := container.CheckImagesSource(ctx, instance.Spec.Template.Spec, a.Config.BackupRegistry)

	if imagesToBackupExist != false {
		return admission.Denied(fmt.Sprintf("There are images to backup %v. Change your image for a backed one in %s before update the resource. ", imagesToBackup, a.Config.BackupRegistry))
	}

	return admission.Allowed(fmt.Sprintf("All the images used by the resource %s/%s on %s are backed.",instance.Kind, instance.Name, instance.Namespace ))
}

func (a *DeploymentImageValidator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
