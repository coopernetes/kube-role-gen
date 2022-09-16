package rolegen

import (
	"flag"
	"fmt"
	"os"
	"log"
	"bytes"
	jsonSer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	client "github.com/coopernetes/kube-role-gen/pkg/k8s"
	gen "github.com/coopernetes/kube-role-gen/pkg/k8s/rbac"
)

func Run() {
	name := flag.String("name", "foo-clusterrole", "Override the name of the ClusterRole "+
		"resource that is generated")
	verbose := flag.Bool("v", false, "Enable verbose logging.")
	json := flag.Bool("json", false, "Generate JSON output. If unset, will default to YAML.")
	jsonPretty := flag.Bool("pretty", false, "Enable human-readable JSON output. This flag is ignored for YAML (always pretty-prints).")
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file. "+
		"If set, this will override the default behavior and "+
		"ignore KUBECONFIG environment variable and/or $HOME/.kube/config file location.")
	flag.Parse()
	d, err := client.SetupDiscoveryClient(*kubeconfig)
	if err != nil {
		fmt.Errorf("Unable to setup client!")
		os.Exit(1)
	}
	cr := gen.GetOrderedResources(*d, *name, *verbose)
	options := jsonSer.SerializerOptions{
		Yaml:   !*json,
		Pretty: *jsonPretty,
	}
	serializer := jsonSer.NewSerializerWithOptions(jsonSer.DefaultMetaFactory, nil, nil, options)
	var writer = bytes.NewBufferString("")
	e := serializer.Encode(cr, writer)
	if e != nil {
		log.Printf("Error encountered during YAML encoding, %s", e.Error())
		os.Exit(1)
	}
	fmt.Println(writer.String())

}
