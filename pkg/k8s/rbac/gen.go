package gen

import (
	"strings"
	"log"
	"os"
	"sort"
	"github.com/coopernetes/kube-role-gen/internal/util"
	"github.com/elliotchance/orderedmap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

func GetOrderedResources(d discovery.DiscoveryClient, name string, verbose bool) *rbacv1.ClusterRole {
	_, resourceListArr, err := d.ServerGroupsAndResources()
	if err != nil {
		log.Printf("Error during resource discovery, %s", err.Error())
		os.Exit(1)
	}

	ordered := orderedmap.NewOrderedMap()
	for _, resourceList := range resourceListArr {
		if verbose {
			log.Printf("Group: %s", resourceList.GroupVersion)
		}
		// rbac rules only look at API group names, not name & version
		groupOnly := strings.Split(resourceList.GroupVersion, "/")[0]
		// core API doesn't have a group "name". We set to "core" and replace at the end with a blank string in the rbac policy rule
		if resourceList.GroupVersion == "v1" {
			groupOnly = "core"
		}

		resourceVerbMap := make(map[string][]string)
		for _, resource := range resourceList.APIResources {
			if verbose {
				log.Printf("Resource: %s - Verbs: %s",
					resource.Name,
					resource.Verbs.String())
			}

			verbs := make([]string, 0)
			for _, v := range resource.Verbs {
				verbs = append(verbs, v)
			}
			sort.Strings(verbs)
			verbString := strings.Join(verbs[:], ",")
			if value, ok := resourceVerbMap[verbString]; ok {
				resourceVerbMap[verbString] = append(value, resource.Name)
			} else {
				resourceVerbMap[verbString] = []string{resource.Name}
			}
		}

		for k := range resourceVerbMap {
			var sb strings.Builder
			sb.WriteString(groupOnly)
			sb.WriteString("!")
			sb.WriteString(k)
			if resourceVal, exists := ordered.Get(sb.String()); exists {
				resourceSetMap := make(map[string]bool)
				for _, r := range resourceVal.([]string) {
					resourceSetMap[r] = true
				}
				for _, r := range resourceVerbMap[k] {
					resourceSetMap[r] = true
				}
				resourceSet := set.MakeSet(resourceSetMap)
				ordered.Set(sb.String(), resourceSet)
			} else {
				ordered.Set(sb.String(), resourceVerbMap[k])
			}
		}
	}

	computedPolicyRules := make([]rbacv1.PolicyRule, 0)
	for _, k := range ordered.Keys() {
		splitKey := strings.Split(k.(string), "!")
		if len(splitKey) != 2 {
			log.Fatalf("Unexpected output from API: %s", k)
		}
		splitVerbList := strings.Split(splitKey[1], ",")
		apiGroup := splitKey[0]
		if splitKey[0] == "core" {
			apiGroup = ""
		}

		value, _ := ordered.Get(k)

		newPolicyRule := &rbacv1.PolicyRule{
			APIGroups: []string{apiGroup},
			Verbs:     splitVerbList,
			Resources: value.([]string),
		}
		computedPolicyRules = append(computedPolicyRules, *newPolicyRule)
	}
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: computedPolicyRules,
	}
}
