package controllers
//
//import (
//	"context"
//	"fmt"
//	"github.com/cargaona/kubermatic-challenge/pkg/container"
//	"github.com/go-logr/logr"
//	appsv1 "k8s.io/api/apps/v1"
//	"k8s.io/apimachinery/pkg/api/errors"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//	"sigs.k8s.io/controller-runtime/pkg/reconcile"
//)
//
//type ReconcileDeployment struct {
//	Client client.Client
//	Logger logr.Logger
//}
//
////Do we need this one?
//var _ reconcile.Reconciler = &ReconcileDeployment{}
//
//func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
//	instance := &appsv1.Deployment{}
//
//	// Ignore all pods on kube-system namespace.
//	if request.Namespace == "kube-system" {
//		return reconcile.Result{}, nil
//	}
//
//	err := r.Client.Get(ctx, request.NamespacedName, instance)
//	if errors.IsNotFound(err) {
//		r.Logger.Error(err, "Could not find Daemonset")
//		return reconcile.Result{}, nil
//	}
//	// False if is from dockerhub true if is from backupRegistry
//	if container.CheckImagesSource(ctx, instance) {
//		r.Logger.Info(fmt.Sprintf("Image already backuped for application: %s", instance.Name))
//		return reconcile.Result{}, nil
//	}
//
//	image, err := container.CopyImagesToBackUpRegistry(ctx, instance)
//	if err != nil {
//		r.Logger.Error(err, "Failed to copy the original image to the backup registry")
//		return reconcile.Result{}, err
//	}
//
//	err = container.UpdateImagesFromDaemonSet(ctx, instance, image)
//	if err != nil {
//		r.Logger.Error(err, "Failed to update the image of the given resource")
//		return reconcile.Result{}, err
//	}
//
//	err = container.ValidateRedeployedDaemonset(ctx, instance, image)
//	if err != nil {
//		r.Logger.Error(err, "The validation was not successful.")
//		return reconcile.Result{}, err
//	}
//
//	r.Logger.Info("Succesfully backuped and redeployed")
//	return reconcile.Result{}, nil
//}
//