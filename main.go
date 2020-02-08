package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sJson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
		log.Printf("ERROR! Unable to build a valid Kubernetes config, %s", err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error during Kubernetes client initialization, %s", err.Error())
		return
	}

	apiResourceListArray, err := clientset.Discovery().ServerResources()
	if err != nil {
		log.Printf("Error during server resource discovery, %s", err.Error())
		return
	}

	computedPolicyRules := make([]rbacv1.PolicyRule, 0)

	for _, apiResourceList := range apiResourceListArray {

		log.Printf("\n\tAPI Group: %s", apiResourceList.GroupVersion)
		// rbac rules only look at API group names, not name & version
		groupOnly := strings.Split(apiResourceList.GroupVersion, "/")[0]
		// core API doesn't have a group "name". In rbac policy rules, its a blank string
		if groupOnly == "v1" {
			groupOnly = ""
		}

		resourceList := make([]string, 0)
		uniqueVerbs := make(map[string]bool)
		for _, apiResource := range apiResourceList.APIResources {
			log.Printf("\n\tGroup Name: %s\n\tResource: %s\n\tVerbs: %s\n",
				groupOnly,
				apiResource.Name,
				apiResource.Verbs.String())

			resourceList = append(resourceList, apiResource.Name)
			for _, verb := range apiResource.Verbs {
				uniqueVerbs[verb] = true
			}
		}

		verbList := make([]string, len(uniqueVerbs))
		i := 0
		for k := range uniqueVerbs {
			verbList[i] = k
			i++
		}

		newPolicyRule := &rbacv1.PolicyRule{
			APIGroups: []string{groupOnly},
			Verbs:     verbList,
			Resources: resourceList,
		}

		computedPolicyRules = append(computedPolicyRules, *newPolicyRule)

	}

	completeRbac := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo-role",
		},
		Rules: computedPolicyRules,
	}

	serializer := k8sJson.NewYAMLSerializer(k8sJson.DefaultMetaFactory, nil, nil)
	var writer = bytes.NewBufferString("")
	e := serializer.Encode(completeRbac, writer)
	if e != nil {
		log.Printf("Error encountered during YAML encoding, %s", e.Error())
		return
	}
	fmt.Println(writer.String())
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
