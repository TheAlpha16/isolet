package instance

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type K8sDeployment struct {
	Client     *K8sClient
	Deployment *appsv1.Deployment
}

func NewK8sDeployment(client *K8sClient, obj *unstructured.Unstructured) (*K8sDeployment, error) {
	if obj == nil {
		return &K8sDeployment{
			Client: client,
		}, nil
	}

	deployment := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, deployment)
	if err != nil {
		return nil, err
	}

	return &K8sDeployment{
		Client:     client,
		Deployment: deployment,
	}, nil
}

func (d *K8sDeployment) Create(ctx context.Context) error {
	if d.Deployment == nil {
		return errors.New("invalid deployment manifest")
	}

	_, err := d.Client.ClientSet.AppsV1().
		Deployments(d.Deployment.Namespace).
		Create(ctx, d.Deployment, metav1.CreateOptions{})
	return err
}

func (d *K8sDeployment) Get(ctx context.Context, name, namespace string) (*appsv1.Deployment, error) {
	return d.Client.ClientSet.AppsV1().
		Deployments(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (d *K8sDeployment) Update(ctx context.Context) error {
	if d.Deployment == nil {
		return errors.New("invalid deployment manifest")
	}

	_, err := d.Client.ClientSet.AppsV1().
		Deployments(d.Deployment.Namespace).
		Update(ctx, d.Deployment, metav1.UpdateOptions{})
	return err
}

func (d *K8sDeployment) Delete(ctx context.Context, name, namespace string) error {
	return d.Client.ClientSet.AppsV1().
		Deployments(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}
