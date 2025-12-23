package main

import (
    "github.com/BurntSushi/toml"
    "os"
)

const CONFIGS_FOLDER string = "tpsg_configs"
const CONFIG_FILE string = "config.toml"
const USERS_CONFIG_FILE string = "users.json"

var TConfig *GKVS = NewGKVS()

func ReadConfig(configPath string) (TConfigTOML, error) {
    var config TConfigTOML

    data, err := os.ReadFile(configPath)
    if err != nil {
        return config, err
    }

    err = toml.Unmarshal(data, &config)
    if err != nil {
        return config, err
    }

    return config, nil
}
