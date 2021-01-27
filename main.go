package main

import (
	"github.com/cargaona/kubermatic-challenge/pkg/controllers"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	//Manager setup
	entryLog.Info("Setting Up Manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "Unable to set up Controller-Manager.")
		os.Exit(1)
	}

	// Setup a new controller to reconcile deployments
	//entryLog.Info("Setting Up Deployment Controller")
	//deploymentReconciler := &controllers.ReconcileDeployment{Client: mgr.GetClient(), Logger: entryLog}
	//deploymentController, err := controller.New("deployment-image-cloner", mgr,
	//	controller.Options{Reconciler: deploymentReconciler,
	//		MaxConcurrentReconciles: 5})

	//if err != nil {
	//	entryLog.Error(err, "Unable to set up Deployment Controller.")
	//	os.Exit(1)
	//}

	entryLog.Info("Setting Up Daemonset Controller")
	daemonsetReconciler := &controllers.ReconcileDaemonset{Client: mgr.GetClient(), Logger: entryLog}
	daemonsetController, err := controller.New("daemonset-image-cloner", mgr,
		controller.Options{Reconciler: daemonsetReconciler,
			MaxConcurrentReconciles: 5})

	if err != nil {
		entryLog.Error(err, "Unable to set up Daemonset controller.")
		os.Exit(1)
	}

	//if err := deploymentController.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
	//	entryLog.Error(err, "Unable to set up watcher for deployments.")
	//	os.Exit(1)
	//}

	if err := daemonsetController.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "Unable to set up watcher for daemonsets.")
		os.Exit(1)
	}
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "Unable to run manager")
		os.Exit(0)
	}
}
