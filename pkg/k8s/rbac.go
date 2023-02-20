// Package k8s provides functions to create Kubernetes RBAC roles objects based
// on discovered API resources. It also provides utility functions to setup a
// discovery client for a provided kubeconfig and obtain the list of discovered
// resources for use in this package.
package k8s

import (
	"github.com/elliotchance/orderedmap/v2"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"sort"
	"strings"
)

// CreateGranularRole creates a ClusterRole where each rules entry contains only the specific combination of API group
// and supported verbs for each resource. Resources with matching verbs are grouped together in a single PolicyRule.
// This differs from other implementations such as `kubectl create clusterrole` which will group together resources
// with verbs that are not applicable or supported.
//
// All PolicyRules in the ClusterRole this function returns represents a "matrix" of all resources available on the API
// and contains only the list of the supported verbs that resource handles.
func CreateGranularRole(apiResourceList []*metav1.APIResourceList, name string, verbose bool) *rbacv1.ClusterRole {
	oMap := orderedmap.NewOrderedMap[string, map[string][]string]()
	for _, resourceList := range apiResourceList {
		if verbose {
			log.Printf("Group %s contains %d resources", resourceList.GroupVersion, len(resourceList.APIResources))
		}
		groupName := extractGroupFromVersion(resourceList.GroupVersion)
		if slices.Contains(oMap.Keys(), groupName) {
			left, _ := oMap.Get(groupName)
			right := convertToVerbMap(resourceList.APIResources, verbose)
			oMap.Set(groupName, mergeVerbMaps(left, right))
		} else {
			oMap.Set(groupName, convertToVerbMap(resourceList.APIResources, verbose))
		}
	}
	policyRules := policyRuleByOrderedMap(*oMap)
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: policyRules,
	}
}

func extractGroupFromVersion(groupVersion string) string {
	if groupVersion == "v1" {
		return ""
	}
	return strings.Split(groupVersion, "/")[0]
}

func convertToVerbMap(resList []metav1.APIResource, verbose bool) map[string][]string {
	verbMap := make(map[string][]string)
	for _, res := range resList {
		if verbose {
			log.Printf("Resource: %s - Verbs: %s",
				res.Name,
				res.Verbs.String())
		}
		verbs := make([]string, len(res.Verbs))
		copy(verbs, res.Verbs)
		sort.Strings(verbs)
		verbKey := strings.Join(verbs, ",")
		if val, ok := verbMap[verbKey]; ok {
			verbMap[verbKey] = append(val, res.Name)
		} else {
			verbMap[verbKey] = []string{res.Name}
		}
	}
	for k := range verbMap {
		sort.Strings(verbMap[k])
	}
	return verbMap
}

func mergeVerbMaps(left map[string][]string, right map[string][]string) map[string][]string {
	merged := make(map[string][]string)
	for k := range left {
		merged[k] = left[k]
	}
	for k := range right {
		var values []string
		if val, ok := merged[k]; ok {
			set := make(map[string]bool)
			for _, v := range right[k] {
				set[v] = true
			}
			for _, v := range val {
				set[v] = true
			}
			toUpdate := mapToSet(set)
			values = make([]string, len(toUpdate))
			copy(values, toUpdate)
		} else {
			values = make([]string, len(right[k]))
			copy(values, right[k])
		}
		sort.Strings(values)
		merged[k] = values
	}
	return merged
}

func policyRuleByOrderedMap(oMap orderedmap.OrderedMap[string, map[string][]string]) []rbacv1.PolicyRule {
	policyRules := make([]rbacv1.PolicyRule, 0)
	for _, group := range oMap.Keys() {
		verbMap, _ := oMap.Get(group)
		for verbStr := range verbMap {
			verbs := strings.Split(verbStr, ",")
			groupName := group
			pr := &rbacv1.PolicyRule{
				APIGroups: []string{groupName},
				Verbs:     verbs,
				Resources: verbMap[verbStr],
			}
			policyRules = append(policyRules, *pr)
		}
	}
	return policyRules
}

func mapToSet(m map[string]bool) []string {
	s := make([]string, 0)
	for k := range m {
		s = append(s, k)
	}
	return s
}
