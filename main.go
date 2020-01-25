package main

import (
	"bytes"
	"flag"
	"fmt"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sJson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	completeRbac := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo-role",
		},
	}
	apiResourceListArray, err := clientset.Discovery().ServerResources()
	if err != nil {
		panic(err.Error())
	}
	for _, apiResourceList := range apiResourceListArray {
		resources := []string{}

		for _, apiResource := range apiResourceList.APIResources {
			fmt.Printf("Group Name: %s\tResource: %s\n", apiResourceList.GroupVersion, apiResource.Name)
			resources = append(resources, apiResource.Name)
			verbs := apiResource.Verbs.String()
			fmt.Printf("Verbs: %s", verbs)
			// s := []string{verbs, apiResourceList.Kind}

		}
		// rbacRules[createKey(s)] = "myValue"
	}
	serializer := k8sJson.NewSerializer(k8sJson.DefaultMetaFactory, nil, nil, true)
	var writer = bytes.NewBufferString("")
	e := serializer.Encode(completeRbac, writer)
	if e != nil {
		panic(e.Error())
	}
	fmt.Printf("Complete RBAC object: %s", writer.String())
}

func createKey(s []string) string {
	return fmt.Sprintf("%q", s)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
