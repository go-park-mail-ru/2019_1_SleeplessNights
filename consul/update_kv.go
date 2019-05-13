package main

import (
	"encoding/json"
	"flag"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/console"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/icza/dyno"
	"github.com/opencopilot/consulkvjson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
)

var consulAddr = flag.String("addr", "127.0.0.1:8500", "Pass consul client address to script")

func main() {
	console.Title("Uploading project configuration to consul...")
	flag.Parse()
	ipRaw, _, err := net.SplitHostPort(*consulAddr)
	if err != nil {
		console.Error("Invalid addr %s : %s", *consulAddr, err.Error())
		return
	}
	ip := net.ParseIP(ipRaw)
	if ip.To4() == nil {
		console.Error("Invalid ipv4 addr %s", *consulAddr)
		return
	}

	yamlKV, err := ioutil.ReadFile("kv.yaml")
	if err != nil {
		console.Error("Can't read kv.yaml: " + err.Error())
		return
	}

	var kvData interface{}
	err = yaml.Unmarshal(yamlKV, &kvData)
	if err != nil {
		console.Error("Unable to unmarshal kv.yaml data, err: %s", err.Error())
		return
	}

	kvData = dyno.ConvertMapI2MapS(kvData)
	jsonKV, err := json.Marshal(kvData)
	if err != nil {
		console.Error("Unable to convert kv.yaml data to json, err: %s", err.Error())
		return
	}

	kvs, err := consulkvjson.ToKVs(jsonKV)
	if err != nil {
		console.Error("Can't convert json converted from kv.yaml to consul KV's: " + err.Error())
		return
	}

	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err := consulapi.NewClient(config)
	if err != nil {
		console.Error("Unable to connect consul on %s", *consulAddr)
		return
	}

	for _, kv := range kvs {
		_, err = consul.KV().Put(&consulapi.KVPair{
			Key:   kv.Key,
			Value: []byte(kv.Value),
		}, nil)
		if err != nil {
			console.Error("Error while putting KVPair to consul key: %s, value: %s, error: %s", kv.Key, kv.Value, err.Error())
			return
		}
	}

	console.Success("All kv pairs successfully put to consul")
}
