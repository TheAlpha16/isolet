package instance

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type IK8sResource interface {
	Create(ctx context.Context, obj runtime.Object) error
	Get(ctx context.Context, name, namespace, kind string) (runtime.Object, error)
	Update(ctx context.Context, obj runtime.Object) error
	Delete(ctx context.Context, name, namespace, kind string) error
}

type K8sResource struct {
	Client *K8sClient
}

func NewK8sResource(client *K8sClient) *K8sResource {
	return &K8sResource{
		Client: client,
	}
}

func (r *K8sResource) Create(ctx context.Context, obj *unstructured.Unstructured) error {
	gvk := obj.GroupVersionKind()
	restMapping, err := r.Client.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	_, err = r.Client.DynamicClient.
		Resource(restMapping.Resource).
		Namespace(obj.GetNamespace()).
		Create(ctx, obj, metav1.CreateOptions{})
	return err
}

func (r *K8sResource) Get(ctx context.Context, name, namespace, kind string) (*unstructured.Unstructured, error) {
	return r.Client.DynamicClient.
		Resource((&unstructured.Unstructured{}).GroupVersionKind().GroupVersion().WithResource(kind)).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (r *K8sResource) Update(ctx context.Context, obj *unstructured.Unstructured) error {
	gvk := obj.GroupVersionKind()
	restMapping, err := r.Client.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	_, err = r.Client.DynamicClient.
		Resource(restMapping.Resource).
		Namespace(obj.GetNamespace()).
		Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func (r *K8sResource) Delete(ctx context.Context, name, namespace, kind string) error {
	return r.Client.DynamicClient.
		Resource((&unstructured.Unstructured{}).GroupVersionKind().GroupVersion().WithResource(kind)).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}
