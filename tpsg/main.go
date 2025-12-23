package main

import (
	"fmt"
	"os"
)

func main() {
	var user_folder string = os.Getenv("HOME")
	var configs_folder_path string = fmt.Sprintf("%s/%s", user_folder, CONFIGS_FOLDER)
	var config_fullpath string = fmt.Sprintf("%s/%s", configs_folder_path, CONFIG_FILE)
	var users_config_fullpath string = fmt.Sprintf("%s/%s", configs_folder_path, USERS_CONFIG_FILE)

	TConfig.Set("user_folder", NewGKVSString(user_folder))
	TConfig.Set("configs_folder_path", NewGKVSString(configs_folder_path))
	TConfig.Set("config_fullpath", NewGKVSString(config_fullpath))
	TConfig.Set("users_config_fullpath", NewGKVSString(users_config_fullpath))

	config, err := ReadConfig(config_fullpath)
	if err != nil {
		LogError(fmt.Sprintf("Failed to read config: %s", err.Error()))
	} else {
		TConfig.Set("config", NewGKVSTConfigTOML(config))
		LogEvent(fmt.Sprintf("Config loaded."))
	}

	err = ReadUsersConfig(users_config_fullpath)
	if err != nil {
		LogError(fmt.Sprintf("Failed to read users config: %s", err.Error()))
	} else {
		LogEvent(fmt.Sprintf("Users config loaded."))
	}

	user_folder_r := TConfig.Get("user_folder").String
	configs_folder_path_r := TConfig.Get("configs_folder_path").String
	config_fullpath_r := TConfig.Get("config_fullpath").String
	users_config_fullpath_r := TConfig.Get("users_config_fullpath").String

	LogInfo(fmt.Sprintf("user_folder: %s", user_folder_r))
	LogInfo(fmt.Sprintf("configs_folder_path: %s", configs_folder_path_r))
	LogInfo(fmt.Sprintf("config_fullpath: %s", config_fullpath_r))
	LogInfo(fmt.Sprintf("users_config_fullpath: %s", users_config_fullpath_r))

	config_r := TConfig.Get("config").TConfigTOML
	LogInfo(fmt.Sprintf("Config TCP: %d", config_r.TCP))
	LogInfo(fmt.Sprintf("Config WS: %d", config_r.WS))
}
