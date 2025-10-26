package sdk

import (
	"context"

	isoletv1 "github.com/TheAlpha16/isolet/tide/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type isoletClient struct {
	k8s client.Client
}

func (isoc *isoletClient) CreateInstance(ctx context.Context, inst *isoletv1.Instance) error {
	return isoc.k8s.Create(ctx, inst)
}

func (isoc *isoletClient) GetInstance(ctx context.Context, name, namespace string) (*isoletv1.Instance, error) {
	var inst isoletv1.Instance
	if err := isoc.k8s.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, &inst); err != nil {
		return nil, err
	}
	return &inst, nil
}

func (isoc *isoletClient) ListInstances(ctx context.Context, namespace string) ([]isoletv1.Instance, error) {
	var instList isoletv1.InstanceList
	if err := isoc.k8s.List(ctx, &instList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}
	return instList.Items, nil
}

func (isoc *isoletClient) UpdateInstance(ctx context.Context, inst *isoletv1.Instance) error {
	return isoc.k8s.Update(ctx, inst)
}

func (isoc *isoletClient) DeleteInstance(ctx context.Context, name, namespace string) error {
	inst := &isoletv1.Instance{}
	inst.Name = name
	inst.Namespace = namespace
	return isoc.k8s.Delete(ctx, inst)
}

func New(k8s client.Client) Handler {
	return &isoletClient{k8s: k8s}
}
