package k8s

import (
	"context"
	"fmt"
	"time"

	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"

	healthSrv "control_plane/internal/service/health"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	apiService "control_plane/internal/service/api_service"
)

type K8sOrchestrator struct {
	clientset           *kubernetes.Clientset
	namespace           string
	image               string
	proxyConnectTimeout string
	proxyReadTimeout    string
	proxySendTimeout    string
	log                 logger.Logger
	health              healthSrv.HealthService
	apiServiceService   apiService.APIServiceService
	clientRepo          repository.ClientRepository
}

func NewK8sOrchestrator(
	clientset *kubernetes.Clientset,
	namespace string,
	image string,
	proxyConnectTimeout string,
	proxyReadTimeout string,
	proxySendTimeout string,
	log logger.Logger,
	health healthSrv.HealthService,
	apiServiceService apiService.APIServiceService,
	clientRepo          repository.ClientRepository,
) orchestrator.Orchestrator {
	return &K8sOrchestrator{
		clientset:           clientset,
		namespace:           namespace,
		image:               image,
		proxyConnectTimeout: proxyConnectTimeout,
		proxyReadTimeout:    proxyReadTimeout,
		proxySendTimeout:    proxySendTimeout,
		log:                 log,
		health:              health,
		apiServiceService:   apiServiceService,
		clientRepo:          clientRepo,
	}
}

func (o *K8sOrchestrator) Deploy(
	ctx context.Context,
	client *domain.APIClient,
	config *domain.APIClientConfig,
) error {
	created := false
	deployName := fmt.Sprintf("api-client-%s", client.ID)

	o.log.Info("k8s deploy started",
		"client_id", client.ID,
		"deploy", deployName,
		"namespace", o.namespace,
	)

	selectorLabels := map[string]string{
		"app":       "api-client",
		"client_id": client.ID,
	}

	templateLabels := map[string]string{
		"app":       "api-client",
		"client_id": client.ID,
		"version":   config.Version,
	}

	apiService, err := o.apiServiceService.GetByID(ctx, client.APIServiceID)
	if err != nil {
		o.log.Error("failed to get api service",
			"client_id", client.ID,
			"api_service_id", client.APIServiceID,
			"error", err,
		)
		return err
	}

	baseURL := apiService.BaseURL

	if baseURL == "" {
		return fmt.Errorf("api service base url is empty")
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: o.namespace,
			Labels:    selectorLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: templateLabels,
					Annotations: map[string]string{
						"config-version": config.Version,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            fmt.Sprintf("client-%s", client.ID),
							Image:           o.image,
							ImagePullPolicy: v1.PullNever,
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("128Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("500m"),
									v1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "CLIENT_ID",
									Value: client.ID,
								},
								{
									Name:  "CLIENT_SLUG",
									Value: client.Slug,
								},
								{
									Name:  "BASE_URL",
									Value: baseURL,
								},
								{
									Name:  "TIMEOUT_MS",
									Value: fmt.Sprintf("%d", config.TimeoutMs),
								},
								{
									Name:  "RETRY_COUNT",
									Value: fmt.Sprintf("%d", config.RetryCount),
								},
								{
									Name:  "RETRY_BACKOFF",
									Value: fmt.Sprintf("%d", config.RetryBackoff),
								},
								{
									Name:  "AUTH_TYPE",
									Value: string(config.AuthType),
								},
								{
									Name:  "AUTH_REF",
									Value: config.AuthRef,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = o.clientset.AppsV1().
		Deployments(o.namespace).
		Create(ctx, deploy, metav1.CreateOptions{})

	if err != nil {

		if apierrors.IsAlreadyExists(err) {
			o.log.Info("deployment already exists, updating",
				"client_id", client.ID,
				"deploy", deployName,
			)

			existing, getErr := o.clientset.AppsV1().
				Deployments(o.namespace).
				Get(ctx, deployName, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}

			existing.Spec = deploy.Spec

			_, updateErr := o.clientset.AppsV1().
				Deployments(o.namespace).
				Update(ctx, existing, metav1.UpdateOptions{})

			if updateErr != nil {
				o.log.Error("k8s update failed",
					"error", updateErr,
				)
				return updateErr
			}

			o.log.Info("k8s deployment updated",
				"client_id", client.ID,
			)
			created = true
		} else {
			o.log.Error("k8s deploy failed",
				"client_id", client.ID,
				"error", err,
			)
			return err
		}
	}

	if created {
		o.log.Info("deployment created")
	} else {
		o.log.Info("deployment updated")
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: o.namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: selectorLabels,
			Ports: []v1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				},
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}

	_, err = o.clientset.CoreV1().
		Services(o.namespace).
		Create(ctx, service, metav1.CreateOptions{})

	if err != nil {
		if apierrors.IsAlreadyExists(err) {

			existing, getErr := o.clientset.CoreV1().
				Services(o.namespace).
				Get(ctx, deployName, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}

			existing.Spec = service.Spec

			_, updateErr := o.clientset.CoreV1().
				Services(o.namespace).
				Update(ctx, existing, metav1.UpdateOptions{})

			if updateErr != nil {
				return updateErr
			}

		} else {
			return err
		}
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: o.namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/$2",
				"nginx.ingress.kubernetes.io/use-regex":      "true",

				"nginx.ingress.kubernetes.io/proxy-connect-timeout": o.proxyConnectTimeout,
				"nginx.ingress.kubernetes.io/proxy-read-timeout":    o.proxyReadTimeout,
				"nginx.ingress.kubernetes.io/proxy-send-timeout":    o.proxySendTimeout,
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptr("nginx"),
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: fmt.Sprintf("/api/%s(/|$)(.*)", client.Slug),
									PathType: ptr(networkingv1.PathTypeImplementationSpecific),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: deployName,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = o.clientset.NetworkingV1().
		Ingresses(o.namespace).
		Create(ctx, ingress, metav1.CreateOptions{})

	if err != nil {
		if apierrors.IsAlreadyExists(err) {

			existing, getErr := o.clientset.NetworkingV1().
				Ingresses(o.namespace).
				Get(ctx, deployName, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}

			existing.Spec = ingress.Spec
			existing.Annotations = ingress.Annotations

			_, updateErr := o.clientset.NetworkingV1().
				Ingresses(o.namespace).
				Update(ctx, existing, metav1.UpdateOptions{})

			if updateErr != nil {
				return updateErr
			}

		} else {
			return err
		}
	}

	return nil
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
	if apierrors.IsNotFound(err) {
		o.log.Warn("deployment not found, nothing to restart",
			"client_id", clientID,
		)
		return nil
	}

	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}

	restartAt := time.Now().Format(time.RFC3339)

	deploy.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = restartAt

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
		"restart_at", restartAt,
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

	policy := metav1.DeletePropagationForeground

	err := o.clientset.AppsV1().
		Deployments(o.namespace).
		Delete(ctx, deployName, metav1.DeleteOptions{
			PropagationPolicy: &policy,
		})

	if err != nil && !apierrors.IsNotFound(err) {
		o.log.Error("failed to delete deployment",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	err = o.clientset.CoreV1().
		Services(o.namespace).
		Delete(ctx, deployName, metav1.DeleteOptions{})

	if err != nil && !apierrors.IsNotFound(err) {
		o.log.Error("failed to delete service",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	err = o.clientset.NetworkingV1().
		Ingresses(o.namespace).
		Delete(ctx, deployName, metav1.DeleteOptions{})

	if err != nil && !apierrors.IsNotFound(err) {
		o.log.Error("failed to delete ingress",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	o.log.Info("k8s delete completed",
		"client_id", clientID,
	)

	return nil
}

func (o *K8sOrchestrator) CheckHealth(ctx context.Context, clientID string) {
	client, err := o.clientRepo.GetByID(ctx, clientID)
	if err == nil {
		if client.GetStatus() == domain.ClientStatusDisabled ||
		client.GetStatus() == domain.ClientStatusDeleting {
			o.health.Set(clientID, domain.HealthUnhealthy)
			return
		}
	} else {
		o.log.Warn("failed to get client, fallback to k8s check",
			"client_id", clientID,
			"error", err,
		)
	}
	
	pods, err := o.clientset.CoreV1().
		Pods(o.namespace).
		List(ctx, metav1.ListOptions{
			LabelSelector: "client_id=" + clientID,
		})

	if err != nil {
		o.log.Error("health check failed", "error", err)
		o.health.Set(clientID, domain.HealthUnhealthy)
		return
	}

	if len(pods.Items) == 0 {
		o.log.Warn("no pods found", "client_id", clientID)
		o.health.Set(clientID, domain.HealthUnhealthy)
		return
	}

	readyPods := 0
	pendingPods := 0

	for _, pod := range pods.Items {

		if pod.Status.Phase == v1.PodPending {
			pendingPods++
			continue
		}

		if pod.Status.Phase != v1.PodRunning {
			continue
		}

		isReady := true
		for _, c := range pod.Status.ContainerStatuses {
			if !c.Ready {
				isReady = false
				break
			}
		}

		if isReady {
			readyPods++
		}
	}
	
	var status domain.HealthStatus

	switch {
	case readyPods  > 0:
		status = domain.HealthHealthy
	case pendingPods > 0:
		status = domain.HealthDegraded
	default:
		status = domain.HealthUnhealthy
	}

	o.health.Set(clientID, status)

	o.log.Info("health updated",
		"client_id", clientID,
		"status", status,
		"pods_total", len(pods.Items),
		"ready", readyPods,
		"pending", pendingPods,
	)
}

func int32Ptr(i int32) *int32 {
	return &i
}

func ptr[T any](v T) *T {
	return &v
}
