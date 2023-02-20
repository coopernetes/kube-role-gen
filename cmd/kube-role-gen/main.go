// The main package for the kube-role-gen executable
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/coopernetes/kube-role-gen/pkg/k8s"
	jsonSer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"log"
	"os"
)

const Version = "v0.0.6"

func main() {
	name := flag.String("name", "foo-clusterrole", "Override the name of the ClusterRole "+
		"resource that is generated")
	verbose := flag.Bool("v", false, "Enable verbose logging.")
	json := flag.Bool("json", false, "Generate JSON output. If unset, will default to YAML.")
	jsonPretty := flag.Bool("pretty", false, "Enable human-readable JSON output. This flag is ignored for YAML (always pretty-prints).")
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file. "+
		"If set, this will override the default behavior and "+
		"ignore KUBECONFIG environment variable and/or $HOME/.kube/config file location.")
	printVersion := flag.Bool("version", false, "Print version info")
	flag.Parse()
	if *printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	d, err := k8s.SetupDiscoveryClient(*kubeconfig)
	if err != nil {
		log.Printf("Error during client setup: %s", err.Error())
		os.Exit(1)
	}
	_, list, err := d.ServerGroupsAndResources()
	if err != nil {
		log.Printf("Error during resource discovery: %s", err.Error())
		os.Exit(1)
	}
	cr := k8s.CreateGranularRole(list, *name, *verbose)
	if err != nil {
		log.Printf("Error during role creation, %s", err.Error())
		os.Exit(1)
	}
	options := serializerOptions(*json, *jsonPretty)
	serializer := jsonSer.NewSerializerWithOptions(jsonSer.DefaultMetaFactory, nil, nil, options)
	var writer = bytes.NewBufferString("")
	e := serializer.Encode(cr, writer)
	if e != nil {
		log.Printf("Error encountered during encoding, %s", e.Error())
		os.Exit(1)
	}
	fmt.Println(writer.String())
}

func serializerOptions(json bool, pretty bool) jsonSer.SerializerOptions {
	if json {
		return jsonSer.SerializerOptions{Yaml: false, Pretty: pretty}
	}
	return jsonSer.SerializerOptions{Yaml: true, Pretty: false}
}
