package k8s

import (
	"github.com/elliotchance/orderedmap/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"reflect"
	"testing"
)

func TestExtractGroupForCore(t *testing.T) {
	in := "v1"
	out := extractGroupFromVersion(in)
	if out != "" {
		t.Fatalf("Expected blank string, got %s", out)
	}
}

func TestExtractGroupForGroupVersion(t *testing.T) {
	in := "batch/v1"
	out := extractGroupFromVersion(in)
	if out != "batch" {
		t.Fatalf("Expected core, got %s", out)
	}
}

func TestGatherResources(t *testing.T) {
	r1 := &metav1.APIResource{
		Verbs: []string{
			"get",
			"patch",
		},
		Name: "test",
	}
	r2 := &metav1.APIResource{
		Verbs: []string{
			"get",
			"list",
		},
		Name: "test2",
	}
	rList := []metav1.APIResource{*r1, *r2}

	actual := convertToVerbMap(rList, true)
	expected := map[string][]string{"get,patch": {"test"}, "get,list": {"test2"}}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected value %s does not match %s", expected, actual)
	}
}

func TestPolicyRuleByMap(t *testing.T) {
	o := orderedmap.NewOrderedMap[string, map[string][]string]()
	coreMap := make(map[string][]string)
	coreMap["get,list"] = []string{"pods"}
	coreMap["get"] = []string{"pods/exec"}
	o.Set("", coreMap)
	actual := policyRuleByOrderedMap(*o)
	if len(actual) != 2 {
		t.Fatalf("Expected 2 policy rules, got %d (%s)", len(actual), actual)
	}
}

func TestMergeVerbMapsRightHasMore(t *testing.T) {
	l := map[string][]string{"create,get,list": {"cronjobs"}}
	r := map[string][]string{"create,get,list": {"cronjobs", "jobs"}}
	actual := mergeVerbMaps(l, r)
	if !slices.Equal(actual["create,get,list"], []string{"cronjobs", "jobs"}) {
		t.Fatalf("Expected merged maps and got %s", actual["create,get,list"])
	}
}

func TestMergeVerbMapsLeftHasMore(t *testing.T) {
	l := map[string][]string{"create,get,list": {"cronjobs", "jobs"}}
	r := map[string][]string{"create,get,list": {"cronjobs"}}
	actual := mergeVerbMaps(l, r)
	if !slices.Equal(actual["create,get,list"], []string{"cronjobs", "jobs"}) {
		t.Fatalf("Expected merged maps and got %s", actual["create,get,list"])
	}
}

func TestMergeVerbMapsHandlesBatchApi(t *testing.T) {
	l := map[string][]string{
		"create delete deletecollection get list patch update watch": {"cronjobs", "jobs"},
		"get patch update": {"cronjobs/status", "jobs/status"},
	}
	r := map[string][]string{
		"create delete deletecollection get list patch update watch": {"cronjobs"},
		"get patch update": {"cronjobs/status"},
	}
	actual := mergeVerbMaps(l, r)
	if !reflect.DeepEqual(l, actual) {
		t.Fatalf("Expected merged map to equal %s, got %s", l, actual)
	}
}
