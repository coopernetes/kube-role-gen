package k8s

import (
	"log"
	"os"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

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
