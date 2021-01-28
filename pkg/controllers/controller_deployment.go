package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/cargaona/kubermatic-challenge/pkg/configuration"
	"github.com/cargaona/kubermatic-challenge/pkg/container"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDeployment struct {
	Client client.Client
	Logger logr.Logger
	Config   configuration.Config
}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	instance := &appsv1.Deployment{}

	// Ignore all pods on kube-system namespace.
	if request.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}

	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		r.Logger.Error(err, "Could not find Deployment")
		return reconcile.Result{}, nil
	}

	imageToBackupExist, images := container.CheckImagesSource(ctx, instance.Spec.Template.Spec, r.Config.BackupRegistry)
	if !imageToBackupExist {
		r.Logger.Info(fmt.Sprintf("Image already backed for application: %s", instance.Name))
		return reconcile.Result{}, nil
	}

	newImages, err := container.CopyImagesToBackUpRegistry(ctx, images, r.Config.BackupRegistry)
	if err != nil {
		r.Logger.Error(err, "Failed to copy the original image to the backup registry")
		return reconcile.Result{}, err
	}

	// Change the images on for every container within the pod.
	for key, value := range newImages {
		instance.Spec.Template.Spec.Containers[key].Image = value
	}

	// Apply the changes on the live object.
	err = r.Client.Update(ctx, instance)
	if err != nil {
		r.Logger.Error(err, "Error updating the pod")
		return reconcile.Result{}, err
	}

	//TODO build ValidateRedeployedDeployment
	//err = container.ValidateRedeployedDeployment(ctx, instance, image)
	//if err != nil {
	//	r.Logger.Error(err, "The validation was not successful.")
	//	return reconcile.Result{}, err
	//}

	r.Logger.Info(fmt.Sprintf("Reconcile completed for Deployment: %s on Namespace: %s", instance.Name, request.NamespacedName))
	return reconcile.Result{}, nil
}
