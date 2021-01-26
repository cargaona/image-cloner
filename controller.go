package main

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconcilePod struct {
	client client.Client
}

// var _ reconcile.Reconciler = &reconcilePod{}

func (r *reconcilePod) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error){
	return reconcile.Result{}, nil
}