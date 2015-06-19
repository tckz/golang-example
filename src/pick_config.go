package main

import (
	"./util"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"strconv"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s /path/to/yaml [key...]\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func pickup(source map[interface{}]interface{}, keys []string) (foundValue interface{}, found bool) {

	var target interface{} = source
	for _, key := range keys {
		switch targetAsType := target.(type) {
		case map[interface{}]interface{}:
			v, ok := targetAsType[key]
			if !ok {
				return
			}
			target = v
		case []interface{}:
			i, err := strconv.Atoi(key)
			if err != nil || i < 0 {
				log.Fatalf("key=%s is not valid number.", key)
			}
			if len(targetAsType) <= i {
				return
			}
			target = targetAsType[i]
		default:
			return
		}
	}

	return target, true
}

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		usage()
		log.Fatalf("*** yaml must be specified.")
	}

	fnYaml := os.Args[1]
	baseYaml, err := util.LoadYaml(fnYaml)
	if err != nil {
		log.Fatalf("*** Failed to load %s: %v\n", fnYaml, err)
	}

	const ARGV_INDEX_KEY_BEGINS = 2
	foundValue, found := pickup(baseYaml, os.Args[ARGV_INDEX_KEY_BEGINS:])
	if !found {
		return
	}

	body, err := yaml.Marshal(&foundValue)
	if err != nil {
		log.Fatalf("*** Failed marshal: %v\n", err)
	}
	fmt.Printf("%s", string(body))
}
