package main

import (
	"os"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

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
		entryLog.Error(err, "Unable to set up controller-manager.")
		os.Exit(1)
	}

	// Setup a new controller to reconcile deployments
	entryLog.Info("Setting Up Controller")
	reconciler := &reconcilePod{client: mgr.GetClient()}
	c, err := controller.New("image-cloner", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		entryLog.Error(err, "Unable to set up reconciler.")
		os.Exit(1)
	}

	// watch deployments and daemon sets
	if err := c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "Unable to set up watcher for deployments.")
		os.Exit(1)
	}

	if err := c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "Unable to set up watcher for daemonsets")
		os.Exit(1)
	}
}
