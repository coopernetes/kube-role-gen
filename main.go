package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"path/filepath"
	"strings"
	"github.com/elliotchance/orderedmap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sJson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var roleNameArg string
	flag.StringVar(&roleNameArg, "name", "foo-clusterrole", "Override the name of the ClusterRole resource that is generated")

	var enableVerboseLogging bool
	flag.BoolVar(&enableVerboseLogging, "v", false, "Enable verbose logging")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("ERROR! Unable to build a valid Kubernetes config, %s", err.Error())
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error during Kubernetes client initialization, %s", err.Error())
		os.Exit(1)
	}

	apiResourceListArray, err := clientset.Discovery().ServerResources()
	if err != nil {
		log.Printf("Error during server resource discovery, %s", err.Error())
		os.Exit(1)
	}

	resourcesByGroupAndVerb := orderedmap.NewOrderedMap()
	for _, apiResourceList  := range apiResourceListArray {
		if enableVerboseLogging {
			log.Printf("Group: %s", apiResourceList.GroupVersion)
		}
		// rbac rules only look at API group names, not name & version
		groupOnly := strings.Split(apiResourceList.GroupVersion, "/")[0]
		// core API doesn't have a group "name". We set to "core" and replace at the end with a blank string in the rbac policy rule
		if apiResourceList.GroupVersion == "v1" {
			groupOnly = "core"
		}

		resourcesByVerb := make(map[string][]string)
		for _, apiResource := range apiResourceList.APIResources {
			if enableVerboseLogging {
				log.Printf("Resource: %s - Verbs: %s",
					apiResource.Name,
					apiResource.Verbs.String())
			}

			verbList := make([]string, 0)
			for _, verb := range apiResource.Verbs {
				verbList = append(verbList, verb)
			}
			sort.Strings(verbList)
			verbString := strings.Join(verbList[:], ",")
			if value,ok := resourcesByVerb[verbString]; ok {
				resourcesByVerb[verbString] = append(value, apiResource.Name)
			} else {
				resourcesByVerb[verbString] = []string {apiResource.Name}
			}
		}

		for k := range resourcesByVerb {
			var sb strings.Builder
			sb.WriteString(groupOnly)
			sb.WriteString("!")
			sb.WriteString(k)
			if resourceVal,exists := resourcesByGroupAndVerb.Get(sb.String()); exists {
				resourceSetMap := make(map[string]bool);
				for _,r := range resourceVal.([]string) {
					resourceSetMap[r] = true
				}
				for _,r := range resourcesByVerb[k] {
					resourceSetMap[r] = true
				}
				resourceSet := mapSetToList(resourceSetMap)
				resourcesByGroupAndVerb.Set(sb.String(), resourceSet)
			} else {
				resourcesByGroupAndVerb.Set(sb.String(), resourcesByVerb[k])
			}
		}
	}

	computedPolicyRules := make([]rbacv1.PolicyRule, 0)
	for _, k := range resourcesByGroupAndVerb.Keys() {
		splitKey := strings.Split(k.(string), "!")
		if len(splitKey) != 2 {
			log.Fatalf("Unexpected output from API: %s", k)
		}
		splitVerbList := strings.Split(splitKey[1], ",")
		apiGroup := splitKey[0]
		if splitKey[0] == "core" {
			apiGroup = ""
		}

		value, _ := resourcesByGroupAndVerb.Get(k)

		newPolicyRule := &rbacv1.PolicyRule{
			APIGroups: []string{apiGroup},
			Verbs:     splitVerbList,
			Resources: value.([]string),
		}
		computedPolicyRules = append(computedPolicyRules, *newPolicyRule)
	}
	completeRbac := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: roleNameArg,
		},
		Rules: computedPolicyRules,
	}

	serializer := k8sJson.NewYAMLSerializer(k8sJson.DefaultMetaFactory, nil, nil)
	var writer = bytes.NewBufferString("")
	e := serializer.Encode(completeRbac, writer)
	if e != nil {
		log.Printf("Error encountered during YAML encoding, %s", e.Error())
		os.Exit(1)
	}
	fmt.Println(writer.String())
}

func mapSetToList(initialMap map[string]bool) []string {
	list := make([]string, len(initialMap))
	i := 0
	for k := range initialMap {
		list[i] = k
		i++
	}
	return list
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
