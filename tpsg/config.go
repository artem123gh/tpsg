package tpsg

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
)

const CONFIGS_FOLDER string = "tpsg_configs"
const CONFIG_FILE string = "config.toml"
const USERS_CONFIG_FILE string = "users.json"

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

func ReadUsersConfig(usersConfigPath string) error {
	data, err := os.ReadFile(usersConfigPath)
	if err != nil {
		return err
	}

	// Map to hold the parsed JSON: username -> {password}
	var usersData map[string]struct {
		Password string `json:"password"`
	}

	err = json.Unmarshal(data, &usersData)
	if err != nil {
		return err
	}

	// Store each user in TUsers GKVS
	for username, userData := range usersData {
		userCreds := TUserCreds{
			Username: username,
			Password: userData.Password,
		}
		TUsers.Set(username, NewGKVSTUserCreds(userCreds))
	}

	return nil
}
