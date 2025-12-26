package main

import (
    "fmt"
    "os"
    "tpsg"
)

func main() {
    var user_folder string = os.Getenv("HOME")
    var configs_folder_path string = fmt.Sprintf("%s/%s", user_folder, tpsg.CONFIGS_FOLDER)
    var config_fullpath string = fmt.Sprintf("%s/%s", configs_folder_path, tpsg.CONFIG_FILE)
    var users_config_fullpath string = fmt.Sprintf("%s/%s", configs_folder_path, tpsg.USERS_CONFIG_FILE)

    tpsg.TConfig.Set("user_folder", tpsg.NewGKVSString(user_folder))
    tpsg.TConfig.Set("configs_folder_path", tpsg.NewGKVSString(configs_folder_path))
    tpsg.TConfig.Set("config_fullpath", tpsg.NewGKVSString(config_fullpath))
    tpsg.TConfig.Set("users_config_fullpath", tpsg.NewGKVSString(users_config_fullpath))

    config, err := tpsg.ReadConfig(config_fullpath)
    if err != nil {
        tpsg.LogError(fmt.Sprintf("Failed to read config: %s", err.Error()))
    } else {
        tpsg.TConfig.Set("config", tpsg.NewGKVSTConfigTOML(config))
        tpsg.LogEvent(fmt.Sprintf("Config loaded."))
    }

    err = tpsg.ReadUsersConfig(users_config_fullpath)
    if err != nil {
        tpsg.LogError(fmt.Sprintf("Failed to read users config: %s", err.Error()))
    } else {
        tpsg.LogEvent(fmt.Sprintf("Users config loaded."))
    }

    user_folder_r := tpsg.TConfig.Get("user_folder").String
    configs_folder_path_r := tpsg.TConfig.Get("configs_folder_path").String
    config_fullpath_r := tpsg.TConfig.Get("config_fullpath").String
    users_config_fullpath_r := tpsg.TConfig.Get("users_config_fullpath").String

    tpsg.LogInfo(fmt.Sprintf("user_folder: %s", user_folder_r))
    tpsg.LogInfo(fmt.Sprintf("configs_folder_path: %s", configs_folder_path_r))
    tpsg.LogInfo(fmt.Sprintf("config_fullpath: %s", config_fullpath_r))
    tpsg.LogInfo(fmt.Sprintf("users_config_fullpath: %s", users_config_fullpath_r))

    config_r := tpsg.TConfig.Get("config").TConfigTOML
    tpsg.LogInfo(fmt.Sprintf("Config TCP: %d", config_r.TCP))
    tpsg.LogInfo(fmt.Sprintf("Config WS: %d", config_r.WS))

    // Check for testseqs command line argument
    if len(os.Args) > 1 && os.Args[1] == "testseqs" {
        // Run test sequences
        tpsg.LogEvent("Running test sequences...")
        tpsg.TestSeqs()
    } else {
        // Normal workflow
        // Run TCP server
        tpsg.RunTCPServer(config_r.TCP)

        // Run WebSocket server
        tpsg.RunWSServer(config_r.WS)
    }

    // Keep the program running
    select {}
}
