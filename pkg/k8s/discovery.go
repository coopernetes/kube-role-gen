package k8s

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// SetupDiscoveryClient will create a new DiscoveryClient. When the kubeconfig arg is unset, the
// client setup uses the usual default behaviour to load either from KUBECONFIG environment variable
// or the default location (usually $HOME/.kube/config). This is provided via client-go package via
// clientcmd.NewDefaultClientConfigLoadingRules.
//
// If kubeconfig string is non-empty, the client will attempt to load the configuration using this value
// by setting the ExplicitPath field on clientcmd.ClientConfigLoadingRules to override the default
// loading rules.
func SetupDiscoveryClient(kubeconfig string) (*discovery.DiscoveryClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loadingRules.ExplicitPath = kubeconfig
	}

	clientConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	restConfig, err := clientConf.ClientConfig()
	if err != nil {
		log.Printf("Error during Kubernetes client initialization, %s", err.Error())
		os.Exit(1)
	}
	dClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		log.Printf("Error during discovery client setup, %s", err.Error())
		os.Exit(1)
	}
	return dClient, nil
}
