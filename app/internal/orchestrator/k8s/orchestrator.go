package k8s

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sOrchestrator struct {
	clientset *kubernetes.Clientset
	namespace string
	log       *slog.Logger
}

func NewK8sOrchestrator(clientset *kubernetes.Clientset, namespace string, log *slog.Logger) orchestrator.Orchestrator {
	return &K8sOrchestrator{
		clientset: clientset,
		namespace: namespace,
		log:       log,
	}
}

func (o *K8sOrchestrator) Deploy(
	ctx context.Context,
	client *domain.APIClient,
	config *domain.APIClientConfig,
) error {

	deployName := fmt.Sprintf("api-client-%s", client.ID)

	o.log.Info("k8s deploy started",
		"client_id", client.ID,
		"deploy", deployName,
		"namespace", o.namespace,
	)

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deployName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": client.ID,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": client.ID,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "client",
							Image: "nginx", // пока заглушка
						},
					},
				},
			},
		},
	}

	_, err := o.clientset.AppsV1().
		Deployments(o.namespace).
		Create(ctx, deploy, metav1.CreateOptions{})

	if err != nil {
		o.log.Error("k8s deploy failed",
			"client_id", client.ID,
			"deploy", deployName,
			"error", err,
		)
		return err
	}

	o.log.Info("k8s deploy created",
		"client_id", client.ID,
		"deploy", deployName,
	)

	return err
}

func (o *K8sOrchestrator) Restart(ctx context.Context, clientID string) error {
	deployName := "api-client-" + clientID

	o.log.Info("k8s restart started",
		"client_id", clientID,
		"deploy", deployName,
		"namespace", o.namespace,
	)

	deploy, err := o.clientset.AppsV1().
		Deployments(o.namespace).
		Get(ctx, deployName, metav1.GetOptions{})
	if err != nil {
		o.log.Error("k8s restart get failed",
			"client_id", clientID,
			"deploy", deployName,
			"error", err,
		)
		return err
	}

	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}

	deploy.Spec.Template.ObjectMeta.Annotations["restartedAt"] = fmt.Sprintf("%d", time.Now().Unix())

	_, err = o.clientset.AppsV1().
		Deployments(o.namespace).
		Update(ctx, deploy, metav1.UpdateOptions{})

	if err != nil {
		o.log.Error("k8s restart update failed",
			"client_id", clientID,
			"deploy", deployName,
			"error", err,
		)
		return err
	}

	o.log.Info("k8s restart completed",
		"client_id", clientID,
		"deploy", deployName,
	)

	return nil
}

func (o *K8sOrchestrator) Delete(ctx context.Context, clientID string) error {
	deployName := "api-client-" + clientID

	o.log.Info("k8s delete started",
		"client_id", clientID,
		"deploy", deployName,
		"namespace", o.namespace,
	)

	err := o.clientset.AppsV1().
		Deployments(o.namespace).
		Delete(ctx, deployName, metav1.DeleteOptions{})

	if err != nil {
		o.log.Error("k8s delete failed",
			"client_id", clientID,
			"deploy", deployName,
			"error", err,
		)
		return err
	}

	o.log.Info("k8s delete completed",
		"client_id", clientID,
		"deploy", deployName,
	)

	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
