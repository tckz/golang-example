package main

import (
	"./util"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s /path/to/base-yaml /path/to/override-yaml...\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func mergeMap(base map[interface{}]interface{},
	override map[interface{}]interface{}) map[interface{}]interface{} {

	var merge func(base map[interface{}]interface{},
		override map[interface{}]interface{}) map[interface{}]interface{}

	merge = func(base map[interface{}]interface{},
		override map[interface{}]interface{}) map[interface{}]interface{} {

		dest := make(map[interface{}]interface{})

		for k, v := range base {
			dest[k] = v
		}

		for k, v := range override {
			baseValue, existInBase := base[k]
			baseValueAsMap, baseValueIsMap := baseValue.(map[interface{}]interface{})
			overrideValueAsMap, overrideValueIsMap := v.(map[interface{}]interface{})
			if existInBase && overrideValueIsMap && baseValueIsMap {
				dest[k] = merge(baseValueAsMap, overrideValueAsMap)
			} else {
				dest[k] = v
			}
		}

		return dest
	}

	return merge(base, override)
}

func main() {
	flag.Parse()
	switch argc := len(os.Args); {
	case argc < 2:
		usage()
		log.Fatalf("*** Base yaml must be specified.")
	case argc < 3:
		usage()
		log.Fatalf("*** Override yaml must be specified.")
	}

	fnBaseYaml := os.Args[1]
	baseYaml, err := util.LoadYaml(fnBaseYaml)
	if err != nil {
		log.Fatalf("*** Failed to load yaml: %v\n", err)
	}

	merged := baseYaml
	for _, fnOverrideYaml := range os.Args[2:] {
		overrideYaml, err := util.LoadYaml(fnOverrideYaml)
		if err != nil {
			log.Fatalf("*** Failed to load yaml: %v\n", err)
		}

		merged = mergeMap(merged, overrideYaml)
	}

	body, err := yaml.Marshal(&merged)
	if err != nil {
		log.Fatalf("*** Failed to marshal load: %v\n", err)
	}

	fmt.Printf("%s", string(body))
}
