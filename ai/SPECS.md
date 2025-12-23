# User specs for the project.

## Project general structure.

1. Folder `ai/` contains some files related to AI assistant.
2. The file `ai/PROMPTS.md` contains some prepared prompts for AI assistant. AI should never edit this file. AI don't even need to read it, because it just contains some prompts which might be not actual for the moment.
3. The file `ai/RULES.md` contains some project coding and handling rules. AI should never edit this file.
4. The file `ai/SPECS.md` contains project specs described by user. AI should never edit this file.
5. The file `ai/SUMMARY.md` contains project summary. This file should be entirely managed by AI. It should be updated when I ask to update summary of the project.
6. Folder `bins/` shoudl compiled binaries of the project.
7. Folder `other/` may temporarily contain some other source file which AI need to read for understanding some tasks. But these file will not be included to repository. So, in summary there shouln't be references to these files.
8. Folder `tpsg/` is the file of golang module for current porject.
9. Script `build_debug.sh` should be responsible to build `tpsg/` module to binary `bins/tpsg_debug` in debug mode.
10. Script `build_release.sh` should be responsible to build `tpsg/` module to binary `bins/tpsg_debug` in release mode.
11. Script `run_console_debug.sh` should be responsible to run `bins/tpsg_debug` in console.
12. Script `run_console_release.sh` should be responsible to run `bins/tpsg_release` in console.

## Logging.

1. Logging functionality should be described in a separate file `tpsg/logging.go`.
2. There should be 3 seprate functions `LogInfo`, `LogEvent` and `LogError`.
3. Each function should take string argument `message`. Then get current timestamp and construct string like `"Event | <current timestamp in format YYYY.MM.DD HH:mm:ss.<milliseconds>> | <message>"`. First word depends on function it should be `Info`, `Event` or `Error` correspondingly. Then string should be printed to console.
4. Function `LogInfo` should normally not be used by AI when writing code - it is for user to put some debug prints, only in case if I tell to use it.
5. Function `LogEvent` should be used whenever you want put some `println` to log some event. If you need to log some constructed event message you should do it like `LogEvent(fmt.Sprintf(<constructed event message>))`.
6. Function `LogError` should be used whenever you want to log some error. If you need to log some constructed error message you should do it like `LogError(fmt.Sprintf(<constructed error message>))`.

## Types.

1. Project specific types should be implemented in `tpsg/types.go`.

## Global key-value storage - GKVS.

1. In app there should be functionality of global key-value storage implemented in `tpsg/gkvs.go`.
2. There should be possibility to create instances of GKVS. Each instans should be something like hash map with `string` keys and values of `GKVSTypes`. Means under each `string` key it should be able to store value of any type listed in `GKVSTypes`.
3. There should be functions which implements API for the set, get and delete values operations in GKVS instance. It should be function which is applied to GKVS instance and take as argument key and value if appliccable. And returns value of operation. Below some pseudocode for it.
```
GKVSInstance.Set(key string, value GKVSTypes)
// Should create new key-value in storage or replace if key already exists. Returns value which was passed.

GKVSInstance.Get(key string)
// Gets value from storage and returns it as result of GKVSTypes. If key doesn't exist should returns as result something like None value from GKVSTypes.

GKVSInstance.Delete(key string)
// Get value from storage, deletes key-value from storage and returns value as result of GKVSTypes. If there was no key-value in storage returns None value from GKVSTypes.
```
4. In pseudocode above I was referring to something like `None` value in `GKVSTypes`. Can you add something like that to `GKVSTypes`?
5. Functionality to work with storage instance should be available from any goroutines and should be thread safe. Probably mutex for read and write should be used. Each instance should have it's own mutex that they can work independenly. Also it should be implemented in the way to perevent deadlocks.

## Configs.

1. Configs functionality should be implemented in file `tpsg/config.go`.
2. This file also can contain some hard coded constands.
3. General config should be represented in `config.toml` which is located in external folder path of which is constructed during app run. It contains some settings like for example `TCP` - TCP listening port for server and `WS` - websocket listening port for server. More setting can be added later. 
4. In `tpsg/config.go` should be implemented function `ReadConfig` which takes TOML config full path. It should read TOML config, parse it and map to type `TConfigTOML`. For parsing TOML package `github.com/BurntSushi/toml` should be used. Obtained `TConfigTOML` value should be returned as result value.
5. In `tpsg/main.go` in `func main` should be sequence to reaf TOML confir and to store it in `TConfig` GKVS.
6. Users config is stored in according to `users_config_fullpath`.
7. User config is JSON file like this.
```
{
    "username1": {
        "password": "password1"
    },
    "username2": {
        "password": "password2"
    },
    "username3": {
        "password": "password3"
    }
}
```
8. In `tpsg/config.go` should be implemented function `ReadUsersConfig` which tales users config full path. It should read `users.json`, parse it and map it to `TUserCreds` type and each user creds object should be stored to the `TUsers` GKVS instance under the key which corresponds to unsername.

## TCP server.

1. Functionality of TCP server should be implemented in file `tpsg/server_tcp.go`.
2. There should be function `RunTCPServer` which takes port argument `unint16`. It should spawn goroutine which whithin which TCP listener will run.
3. For each TCP connection should be spawned another goroutine with function `HandleTCPConnection`. Within this function requests should be processed synchronously until conncetion is closed.
4. For each request should be called function `ProcessTCPRequest` which takes request as argument and gives response as result. Then result should be sent back to client. For now you can put placeholder for this function which will decode request as text and then will send back echo response to client.
5. In `tpsg/main.go` in `func main` in the and of the function TCP server should be run wth port TCP from TConfig.