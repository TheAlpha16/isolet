package instance

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type K8sService struct {
	Client  *K8sClient
	Service *corev1.Service
}

func NewK8sService(client *K8sClient, obj *unstructured.Unstructured) (*K8sService, error) {
	if obj == nil {
		return &K8sService{
			Client: client,
		}, nil
	}

	service := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, service)
	if err != nil {
		return nil, err
	}

	return &K8sService{
		Client:  client,
		Service: service,
	}, nil
}

func (s *K8sService) Create(ctx context.Context) error {
	if s.Service == nil {
		return errors.New("invalid service manifest")
	}

	_, err := s.Client.ClientSet.CoreV1().
		Services(s.Service.Namespace).
		Create(ctx, s.Service, metav1.CreateOptions{})
	return err
}

func (s *K8sService) Get(ctx context.Context, name, namespace string) (*corev1.Service, error) {
	return s.Client.ClientSet.CoreV1().
		Services(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (s *K8sService) Update(ctx context.Context) error {
	if s.Service == nil {
		return errors.New("invalid service manifest")
	}

	_, err := s.Client.ClientSet.CoreV1().
		Services(s.Service.Namespace).
		Update(ctx, s.Service, metav1.UpdateOptions{})
	return err
}

func (s *K8sService) Delete(ctx context.Context, name, namespace string) error {
	return s.Client.ClientSet.CoreV1().
		Services(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}
