package config

import (
	"bytes"
	"encoding/json"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/opencopilot/consulkvjson"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"time"
)

var consulAddr = os.Getenv("CONSUL_ADDR")

var (
	logger                            = log.GetLogger("ViperConfig")
	consul          *consulapi.Client
	consulLastIndex uint64            = 0
	consulKvPrefix                    = ""
)

func init() {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	var err error
	consul, err = consulapi.NewClient(config)
	if err != nil {
		logger.Fatal("Unable to connect to consul on", consulAddr)
	}

	viper.SetConfigType("json")

	err = loadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	go runConfigUpdater()
}

func runConfigUpdater() {
	ticker := time.Tick(10 * time.Second)
	for range ticker {
		err := loadConfig()
		if err != nil {
			logger.Error("Failed to update configuration:", err)
		}
	}
}

func loadConfig()(err error) {
	qo := &consulapi.QueryOptions{
		WaitIndex: consulLastIndex,
	}
	kvPairs, qm, err :=  consul.KV().List(consulKvPrefix, qo)
	if err != nil {
		logger.Error("Unable to load consul config", err)
		return err
	}

	if consulLastIndex == qm.LastIndex {
		//Consul config didn't change since the last load
		return nil
	}

	jsonKVs, err := consulkvjson.ConsulKVsToJSON(kvPairs)
	if err != nil {
		logger.Error("Unable to convert consul KVs to JSON")
		return err
	}

	jsonData, err := json.Marshal(jsonKVs)
	if err != nil {
		logger.Error("Unable to marshal jsonKVs to json", err)
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Viper failed to read jsonData", err)
	}

	consulLastIndex = qm.LastIndex
	logger.Info("Config updated!")
	return nil
}
