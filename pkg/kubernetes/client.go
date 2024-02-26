package kubernetes

import (
	"context"
	"fmt"
	"github.com/xuxant/kbuild/pkg/output/log"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	DefaultClient = getDefaultClientset
)

func getRestClientConfig(kctx, kcfg string) (*restclient.Config, error) {
	log.Entry(context.TODO()).Debug("Getting Client config for kubernetes")

	rawConfig, err := getCurrentConfig()
	if err != nil {
		return nil, err
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(rawConfig, kctx, &clientcmd.ConfigOverrides{
		CurrentContext: kctx,
	}, clientcmd.NewDefaultClientConfigLoadingRules())
	restConfig, err := clientConfig.ClientConfig()
	if kctx == "" && kcfg == "" && clientcmd.IsEmptyConfig(err) {
		log.Entry(context.TODO()).Debug("No kubeContext set and no kubeConfig found")
		return restConfig, fmt.Errorf("error creating REST client config")
	}
	if err != nil {
		return restConfig, fmt.Errorf("error creating REST client config for kubeContext %q: %w", kctx, kcfg)
	}
	return restConfig, nil
}

func GetDefaultRestClientConfig() (*restclient.Config, error) {
	return getRestClientConfig(kubeContext, kubeConfigFile)
}

func getDefaultClientset() (kubernetes.Interface, error) {
	config, err := GetDefaultRestClientConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}
