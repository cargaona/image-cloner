package controllers

import (
	"context"
	"fmt"
	"github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/container"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDaemonset struct {
	Client client.Client
	Logger logr.Logger
	Config configuration.Config
}

func (r *ReconcileDaemonset) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	instance := &appsv1.DaemonSet{}

	// Ignore all pods on configured ignored namespaces.
	if containsString(r.Config.NamespacesToIgnore, request.Namespace) {
		r.Logger.Info(fmt.Sprintf("Ignoring application since it's deployed on the %s namespace", request.Namespace))
		return reconcile.Result{}, nil
	}

	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		r.Logger.Error(err, "Could not find Daemonset")
		return reconcile.Result{}, nil
	}

	imageToBackupExist, images := container.CheckImagesSource(ctx, instance.Spec.Template.Spec, r.Config.BackupRegistry)
	if !imageToBackupExist {
		r.Logger.Info(fmt.Sprintf("Images already backed for application: %s", instance.Name))
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

	// Apply the changes on the live object
	err = r.Client.Update(ctx, instance)
	if err != nil {
		r.Logger.Error(err, "Error updating the pod")
		return reconcile.Result{}, err
	}

	//TODO pass just the instance object to the function
	err = container.ValidateRedeployedDaemonset(ctx, instance.Status.NumberUnavailable, instance.Spec.Template.Spec, newImages)
	if err != nil {
		r.Logger.Error(err, "The validation was not successful.")
		return reconcile.Result{}, err
	}

	r.Logger.Info(fmt.Sprintf("Reconcile completed for Daemonset: %s on Namespace: %s", instance.Name, request.NamespacedName))
	return reconcile.Result{}, nil
}
