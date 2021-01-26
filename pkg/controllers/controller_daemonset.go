package controllers

import (
	"context"
	"github.com/cargaona/kubermatic-challenge/pkg/container"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDaemonset struct {
	Client client.Client
	Logger logr.Logger
}

//Do we need this one?
var _ reconcile.Reconciler = &ReconcileDeployment{}

func (r *ReconcileDaemonset) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	instance := &appsv1.DaemonSet{}

	// Ignore all pods on kube-system namespace.
	if request.Namespace == "kube-system" {
		r.Logger.Info("Ignoring application since it's deployed on the kube-system namespace")
		return reconcile.Result{}, nil
	}

	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		r.Logger.Error(err, "Could not find Daemonset")
		return reconcile.Result{}, nil
	}
	// False if is from dockerhub true if is from backupRegistry
	if container.CheckImageSource(ctx, instance) {
		r.Logger.Info("Image already backuped")
		return reconcile.Result{}, nil
	}

	image, err := container.CopyImageToBackUpRegistry(ctx, instance)
	if err != nil {
		r.Logger.Error(err, "Failed to copy the original image to the backup registry")
		return reconcile.Result{}, err
	}

	err = container.UpdateImageFromResource(ctx, instance, image)
	if err != nil {
		r.Logger.Error(err, "Failed to update the image of the given resource")
		return reconcile.Result{}, err
	}

	err = container.ValidateRedeployedApplication(ctx, instance, image)
	if err != nil {
		r.Logger.Error(err, "The validation was not successful.")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
