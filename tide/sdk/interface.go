package sdk

import (
	"context"

	isoletv1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler interface {
	CreateInstance(ctx context.Context, inst *isoletv1.Instance) error
	GetInstance(ctx context.Context, name, namespace string) (*isoletv1.Instance, error)
	ListInstances(ctx context.Context, opts ...client.ListOption) ([]isoletv1.Instance, error)
	UpdateInstance(ctx context.Context, inst *isoletv1.Instance) error
	DeleteInstance(ctx context.Context, name, namespace string) error
}
