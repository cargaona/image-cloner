package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/container"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDeployment struct {
	Client client.Client
	Logger logr.Logger
	Config configuration.Config
}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	instance := &appsv1.Deployment{}

	// Ignore all pods on configured ignored namespaces.
	if containsString(r.Config.NamespacesToIgnore, request.Namespace) {
		r.Logger.Info(fmt.Sprintf("Ignoring application since it's deployed on the %s namespace", request.Namespace))
		return reconcile.Result{}, nil
	}

	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		r.Logger.Error(err, "")
		return reconcile.Result{}, nil
	}

	imageToBackupExist, images := container.CheckImagesSource(ctx, instance.Spec.Template.Spec, r.Config.BackupRegistry)
	if !imageToBackupExist {
		r.Logger.Info(fmt.Sprintf("Images already backed for %s/%s", instance.Kind, instance.Name))
		return reconcile.Result{}, nil
	}

	newImages, err := container.CopyImagesToBackUpRegistry(ctx, images, r.Config.BackupRegistry)
	if err != nil {
		r.Logger.Error(err, fmt.Sprintf("failed to copy the original of %s/%s image to the backup registry", instance.Kind, instance.Name))
		return reconcile.Result{}, err
	}

	// Change the images on for every container within the pod.
	for key, value := range newImages {
		instance.Spec.Template.Spec.Containers[key].Image = value
	}

	// Apply the changes on the live object.
	err = r.Client.Update(ctx, instance)
	if err != nil {
		r.Logger.Error(err, fmt.Sprintf( "error updating pods for %s/%s", instance.Kind, instance.Name))
		return reconcile.Result{}, err
	}

	r.Logger.Info(fmt.Sprintf("Reconcile completed for %s/%s on: %s", instance.Kind, instance.Name, request.NamespacedName))
	return reconcile.Result{}, nil
}
