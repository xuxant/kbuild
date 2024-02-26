package kubernetes

import (
	"context"
	"github.com/xuxant/kbuild/pkg/options"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sync"
)

var (
	CurrentConfig = getCurrentConfig
)

var (
	kubeConfigOnce      sync.Once
	kubeConfigFile      string
	kubeContext         string
	kubeConfig          clientcmd.ClientConfig
	checkPermissionOnce sync.Once
)

func getCurrentConfig() (clientcmdapi.Config, error) {
	kubeConfigOnce.Do(func() {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = kubeConfigFile
		kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{
			CurrentContext: kubeContext,
		})
	})
	cfg, err := kubeConfig.RawConfig()
	if kubeContext != "" {
		cfg.CurrentContext = kubeContext
	}
	return cfg, err
}

func CheckPermissions() (bool, map[string]string) {
	cfg := options.GetConfig()
	connection := false
	msg := make(map[string]string)
	checkPermissionOnce.Do(func() {
		config, _ := clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			msg["clientSet"] = err.Error()
			return
		}
		//	Check Create Pod
		_, err = clientSet.CoreV1().Pods(cfg.Namespace).Create(context.TODO(), &corev1.Pod{}, metav1.CreateOptions{})

		if errors.IsForbidden(err) {
			msg["CreatePod"] = err.Error()
			return
		}

		logs := clientSet.CoreV1().Pods(cfg.Namespace).GetLogs("pod-name", &corev1.PodLogOptions{})
		_, err = logs.Stream(context.Background())
		if !errors.IsNotFound(err) {
			msg["PodLog"] = err.Error()
			return
		}
		connection = true

	})
	return connection, msg
}
