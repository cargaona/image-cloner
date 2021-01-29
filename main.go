package main

import (
	clonerConfig "github.com/cargaona/image-cloner/pkg/configuration"
	"github.com/cargaona/image-cloner/pkg/controllers"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"os"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New(zap.UseDevMode(true)))
}

func main() {
	entryLog := log.Log.WithName("image-cloner")

	// Load config
	conf, err := clonerConfig.GetConfig()
	if err != nil {
		entryLog.Error(err, "error loading configuration")
		os.Exit(1)
	}

	//Manager setup
	entryLog.Info("Setting Up Manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up controller-manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile deployments
	entryLog.Info("Setting Up Deployment Controller")
	deploymentReconciler := &controllers.ReconcileDeployment{Client: mgr.GetClient(), Logger: entryLog, Config: *conf}
	deploymentController, err := controller.New("deployment-image-cloner", mgr,
		controller.Options{Reconciler: deploymentReconciler,
			MaxConcurrentReconciles: conf.MaxConcurrentReconciles})

	if err != nil {
		entryLog.Error(err, "unable to set up deployment controller")
		os.Exit(1)
	}

	// Setup a new controller to reconcile daemonsets
	entryLog.Info("Setting Up Daemonset Controller")
	daemonsetReconciler := &controllers.ReconcileDaemonset{Client: mgr.GetClient(), Logger: entryLog, Config: *conf}
	daemonsetController, err := controller.New("daemonset-image-cloner", mgr,
		controller.Options{Reconciler: daemonsetReconciler,
			MaxConcurrentReconciles: conf.MaxConcurrentReconciles})

	if err != nil {
		entryLog.Error(err, "unable to set up daemonset controller")
		os.Exit(1)
	}

	// watch resources
	if err := deploymentController.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to set up watcher for deployments")
		os.Exit(1)
	}

	if err := daemonsetController.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to set up watcher for daemonsets")
		os.Exit(1)
	}

	//Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/mutate-v1-pod", &webhook.Admission{Handler: &controllers.PodAnnotator{Client: mgr.GetClient()}})

	// start manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)

	}
}
