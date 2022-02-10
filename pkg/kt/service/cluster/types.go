package cluster

import (
	"context"
	opt "github.com/alibaba/kt-connect/pkg/kt/options"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesInterface kubernetes interface
type KubernetesInterface interface {
	CreatePod(ctx context.Context, metaAndSpec *PodMetaAndSpec) error
	CreateShadowPod(ctx context.Context, metaAndSpec *PodMetaAndSpec, sshcm string) error
	GetPod(ctx context.Context, name string, namespace string) (*coreV1.Pod, error)
	GetPodsByLabel(ctx context.Context, labels map[string]string, namespace string) (*coreV1.PodList, error)
	UpdatePod(ctx context.Context, pod *coreV1.Pod) (*coreV1.Pod, error)
	WaitPodReady(ctx context.Context, name, namespace string, timeoutSec int) (*coreV1.Pod, error)
	WaitPodTerminate(ctx context.Context, name, namespace string) (*coreV1.Pod, error)
	IncreaseRef(ctx context.Context, name ,namespace string) error
	DecreaseRef(ctx context.Context, name, namespace string) (bool, error)
	AddEphemeralContainer(ctx context.Context, containerName, podName string, envs map[string]string) (string, error)
	RemoveEphemeralContainer(ctx context.Context, containerName, podName string, namespace string) error
	ExecInPod(containerName, podName, namespace string, cmd ...string) (string, string, error)
	RemovePod(ctx context.Context, name, namespace string) error

	GetDeployment(ctx context.Context, name string, namespace string) (*appV1.Deployment, error)
	GetDeploymentsByLabel(ctx context.Context, labels map[string]string, namespace string) (*appV1.DeploymentList, error)
	GetAllDeploymentInNamespace(ctx context.Context, namespace string) (*appV1.DeploymentList, error)
	UpdateDeployment(ctx context.Context, deployment *appV1.Deployment) (*appV1.Deployment, error)
	ScaleTo(ctx context.Context, deployment, namespace string, replicas *int32) (err error)

	CreateService(ctx context.Context, metaAndSpec *SvcMetaAndSpec) (*coreV1.Service, error)
	UpdateService(ctx context.Context, svc *coreV1.Service) (*coreV1.Service, error)
	GetService(ctx context.Context, name, namespace string) (*coreV1.Service, error)
	GetServicesBySelector(ctx context.Context, matchLabels map[string]string, namespace string) ([]coreV1.Service, error)
	GetAllServiceInNamespace(ctx context.Context, namespace string) (*coreV1.ServiceList, error)
	GetServicesByLabel(ctx context.Context, labels map[string]string, namespace string) (*coreV1.ServiceList, error)
	RemoveService(ctx context.Context, name, namespace string) (err error)
	WatchService(name, namespace string, fAdd, fDel, fMod func(*coreV1.Service))

	CreateConfigMapWithSshKey(ctx context.Context, labels map[string]string, sshcm string, namespace string,
		generator *util.SSHGenerator) (configMap *coreV1.ConfigMap, err error)
	GetConfigMap(ctx context.Context, name, namespace string) (*coreV1.ConfigMap, error)
	GetConfigMapsByLabel(ctx context.Context, labels map[string]string, namespace string) (*coreV1.ConfigMapList, error)
	RemoveConfigMap(ctx context.Context, name, namespace string) (err error)

	GetAllNamespaces(ctx context.Context) (*coreV1.NamespaceList, error)
	ClusterCidrs(ctx context.Context, namespace string) (cidrs []string, err error)
}

// Kubernetes implements KubernetesInterface
type Kubernetes struct {
	Clientset kubernetes.Interface
}

// Cli the singleton type
var instance *Kubernetes

// Ins get singleton instance
func Ins() *Kubernetes {
	if instance == nil {
		instance = &Kubernetes{
			Clientset: opt.Get().RuntimeStore.Clientset,
		}
	}
	return instance
}