# Project corntron
Running an extended TOML-defined environment with one or more extended TOML-defined configurations.  
This source code is a reference implementation and prototype design of corntron. 
It is licensed under Mulan PSL v2. 

## How to use corntron?
You can integrate corntron into your program by using the `corntron` package.
Or you can use the prebuilt [`corntron` command](#about-prebuilt-corntron) to run a corntron program.

## What is a corn?
A corn is a program that can be run with a set of arguments and a set of environment variables.
- It's defined by one or more extended TOML files.
- Can be defined to run a program.
- Can be defined as a meta-configuration for a corn.

## What is a runtime?
A runtime is a set of environment variables and a set of configuration variables.
- It's defined by an extended TOML file.
- Can be defined to prepare an environment for corn or interactive commandline.
- Can be defined to prepare public tools (a set of tools for public use, eg. compilers, shells, language runtime,etc.) for a corn or interactive commandline.
- Can be used to automatically do some work when you need mirror settings because of country-specific networking issues.

## What is directory structure?
```
corns
- corns\_env
-- corns\_env\*.toml
runtimes
- runtimes\_env
-- runtimes\_env\*.toml
#{OS}-#{ARCH}
- #{OS}-#{ARCH}/corns
- #{OS}-#{ARCH}/runtimes
```
- `corns` directory contains corn environments and the corn environment root for common platform.
- `corns/_env/<corn-name>.toml` contains corn configurations.
- `runtimes` directory contains runtimes and the runtime environment root for common platform.
- `runtimes/_env/<runtime-name>.toml` directory contains runtimes config.
- `#{OS}-#{ARCH}` directory contains current operating system and architecture specific binaries of corns and runtimes,can be override by `[platform_dir]` of `core.toml`. 
- `#{OS}-#{ARCH}/corns` directory contains current operating system and architecture specific corns.
- `#{OS}-#{ARCH}/runtimes` directory contains current operating system and architecture specific runtimes.
> Current `OS` as [`GOOS`](https://github.com/golang/go/tree/master/src/internal/goos) and `ARCH` [`GOARCH`](https://github.com/golang/go/tree/master/src/internal/goarch) are defined by Go programming Language. These will be replaced by LLVM Triple style [`os`](https://github.com/llvm/llvm-project/blob/main/llvm/include/llvm/TargetParser/Triple.h#L205) and [`arch`](https://github.com/llvm/llvm-project/blob/main/llvm/include/llvm/TargetParser/Triple.h#L46) in the future.

## What content in a corntron's main config (as is core.toml) file
```toml
# base_dir defines the path of the corntron's root dir
# default as the dir of the corntron's executable file without link as internal placeholder `${dp0}`
base_dir = ""

# corn_dirname defines the folder name of corn environments and binaries
# default as "runtimes" 
runtime_dirname = ""

# corn_dirname defines the folder name of corn environments and binaries
# default as "corns" 
corn_dirname = ""

# mirror_type defines the mirror type of current corntron use.
# this mirror type is one of the built-in mirror types and `mirror_types` defined.
# built-in mirror types are `cn`,`none`
mirror_type = ""

# mirror_types defines other mirror type than built-in mirror types
# built-in mirror types are `cn`,`none`
mirror_types = ["private","alternate-1"]

# with_corn defines if this corntron base running the corn config
# defaults as true
with_corn = true

# platform_dir defines binaries dir of corn and runtime environments
# default as `<OS>-<Arch>` when OS is `windows` and architecture is `amd64`,the folder will be `windows-amd64`
[platform_dir]
#This is an example when you need to override them
windows-amd64 = "bin_x64"

# profile_dir defines the $HOME or %USERPROFILE% when running by the corntron
# default as `${userprofile}` 
# `${userprofile}` is an internal placeholder as your host's $HOME or %USERPROFILE%
# `${currentdir}` is an internal placeholder as configure root's directory
profile_dir = ""
```

## What content in a corntron environment config file?
All corntron environment config files are extended TOML files.
Their structure is defined like the following TOML file.
```toml
# This is a corntron environment config file.  

# is_common_platform defines whether this environment is for common platform.
is_common_platform = false

# dir_name defines the directory name of a corntron environment. 
# If empty,the dirname will be the same as the config file name.
dir_name = "" 

# `depend_runtimes` is a set of runtimes that current config depends on.
# The value of `depend_runtimes` is a list of string.
# The reference format is `[<registry>:]<runtime_name>[@<version>]`
# The `<registry>` is dir or network address of the registry.
# The `<version>` will be override the `#{corn_name}_spec_version` variable to pass the version to config.
# The depend runtime config can export a set of `var` to use.
# If the depend runtime are not configured, this config will not be executed.
depend_runtimes = [""]

# `vars` is a set of corntron variable.
# It's used to store some key-value pairs that is not exposed to environment variables.
# Can reference a variable or environment variable by `#{var_name}`.
# Can reference a variable and environment variable with built-in functions by `#{var_name:func_name([param])}`.
[vars]
# This is a normal key name
foo="bar"

# This is a key defined with built-in functions to setting the values.
# If the built-in function is execution failed, the value will be the original value.
"bar:rp(oo=qq)"="foo"

# `mirror_vars.<mirror type>` is a set of corntron variable.
# It's used to override the `vars` when corntron's current `mirror_type` matches.
# The content of the each `<mirror_type>' as same as `vars` usage.
[mirror_vars.cn]

# `envs` is a set of corntron environment variables.
# It's used to store some key-value pairs that is exposed to environment variables.
# Can reference a variable and environment variable by `#{var_name}`.
# Can reference a variable and environment variable with built-in functions by `#{var_name:func_name([param])}`.
[envs]

# This is a environment variable.
# The value of this environment variable will reference to a corntron variable or corntron envrionment variable.
# After referencing, the refrenced value will exectute with the built-in functions.
FOO_VAL="#{foo:rp(oo=qq)}"

# This is a set of corntron bootstrap exection commands.
# It's used to prepare an environment for creating corn or runtime.
# The execution order is the same as the order of the list.
# The content structure of the each element is same as the `exec` element of corn config file. 
[[bootstrap_exec]]

# This is a set of corntron configuration exection commands.
# It's used to prepare an environment for corn or interactive commandline.
# The execution order is the same as the order of the list.
# The content structure of the each element is same as the `exec` element of corn config file.
[[config_exec]]
```
The corn config file is a corntron environment config file that is used to running a program or command. [Here](https://cnb.cool/viscropst/corntron_env_registry/-/tree/main/corns) is the example.
```toml
# This is a corn config file.

# `meta_only` defines whether this corn config is a meta-corn config.
# meta-corn config is a corn config that not executing the commands in `bootstrap_exec`.
# If `meta_only` is true, the corn config is a meta-corn config.
# Otherwise, the corn config is a normal corn config.
meta_only = false

# `import_corns` is a set of corns that this corn imports.
# If the import corns are not running, this corn will not be executed.
# And the `vars` and `envs` of this corn will be merged with the `vars` and `envs` of the import corns.
# The value of `import_corns` is a list of string.
import_corns = ["foo"]

# `exec` is the main execution command of a corn environment.
[exec]
# `exec` is the command name or the file path of a executable file.
# built-in command name will be prefixed with `i-<command>` command name (eg. `i-utar`, `i-cp`).
exec = "echo"

# `args` is the set of arguments that will be passed to the executable file.
# The value of `args` is a list of string.
args = ["Hello World"]

# `platform` is the operating system and architecture of the executable file.
# The value of `platform` is a string.
platform = "windows"

# `work_dir` is the working directory of the executable file.
# The value of `work_dir` is a string.
work_dir = "."

# `with_environ` is the flag that defines whether the environment variables will be passed to the executable file.
# If `with_environ` is true, the environment variables will be passed to the executable file.
# Otherwise, the environment variables will not be passed to the executable file.
# The envrion is the current environment variables when you executing this corn config or intractive commands.
with_environ = false

# `with_no_waiting` is the flag that defines whether the corntron will wait for the execution of the executable file.
# If `with_no_waiting` is true, the corntron will not wait for the execution of the executable file, 
# And execute as a new process.
# Otherwise, the corntron will wait for the execution of the executable file,
# And execute as a child of the corntron process.
with_no_waiting = true

# `is_background` is the flag that defines whether the command will be executed in background.
# The background is a process that will not showing the window.
# If `is_background` is true, the command will be executed in background.
# Otherwise, the command will be executed in foreground.
is_background = false

# `arg_str` is the set of arguments that will be passed to the executable file.
# it will convert a string to a list of string append to `args`.
[exec.arg_str]
# src is the source string that will be converted to a list of string append to `args`.
src="Hello?.Hi?"
# `split_str` is the string that will be used to split the `src` string.
split_str="."
# `split_num` is the number of split string.
split_num=1
# `replaces` is the set of string that will be used to replace the `src` string.
# the first element of `replaces` is the string that will be replaced by the second element of `replaces`.
replaces=["?"," World"]

# `exec.vars` is the set of corntron variables that used for this command only.
# The value of `exec.vars` is a set of key-value pairs.
# The usage of `exec.vars` is the same as `vars` in corntron environment config.
[exec.vars]

# `exec.envs` is the set of environment variables that will be passed to the executable file.
# The value of `exec.envs` is a set of key-value pairs.
# The usage of `exec.envs` is the same as `envs` in corntron environment config.
[exec.envs]

# `exec_by_plat` is the main execution command of a corntron environment for a specific operating system and architecture.
# If this is not defined, the `exec` will be the main execution command.
# Otherwise, the matched `exec_by_plat.<os>-<arch>` will be the main execution command.
# The structure of `exec_by_plat.<os>-<arch>` is the same as `exec` structure in corn evrironment config.
[exec_by_plat.windows-amd64]
```
The runtime environment config file is a corntron environment config file that used to define a runtime environment for a corntron program (a corn config file or interactive commandline).[Here](https://cnb.cool/viscropst/corntron_env_registry/-/tree/main/runtimes) is the example.
```toml
# This is a corntron runtime environment config file.

# `mirror_envs.<mirror-type>` is the set of environment variables that will be used to set mirror settings.
# `<mirror-type>` is the type of mirror settings, current is simplfield name of country code (eg. `cn`).
# You can add more mirror types in core config of corntron (a.k.a. `core.toml`).
[mirror_envs.cn]
FOO_VAL="bar"

# `mirror_exec.<mirror-type>` is the set of execution commands that will be used to set mirror settings.
# `<mirror-type>` is the type of mirror settings, current is simplfield name of country code (eg. `cn`).
# You can add more mirror types in core config of corntron (a.k.a. `core.toml`).
# The structure of `mirror_exec.<mirror-type>` is the same as `exec` structure in corn evrironment config.
[[mirror_exec.cn]]

```

## What are built-in commands?
The built-in commands can be running by `exec`,`config_exec`,`mirror_exec`,`bootstrap_exec`
- `i-cp` is the command that used to copy a file or a directory.
- `i-md` is the command that used to create a directory.
- `i-mv` is the command that used to move a file or a directory.
- `i-rd` is the command that used to a directory and files in directory.
- `i-wgt` is the command that used to download a file from internet.
- `i-ghgt` is the command that used to download a file from github release.
- `i-utar` is the command that used to unarchive a `.tar`,`.tgz`,`.tar.xz`,`.tar.gz` or `.tar.bz2` file.
- `i-uzip` is the command that used to unarchive a `.zip` file.
- `i-wstr` is the command that used to write a string to a file.

## What are built-in functions?
- `rp` is the function that used to replace a string.
- `ospth` is the function that used to convert a path to a operating system path.
- `webreq` is the function that used to send a http request and replace to the response body string.
- `gh-rel-ver` is the function that used to get the latest release version of a github repository.
- `gl-rel-ver` is the function that used to get the latest release version of a gitlab repository.

## About prebuilt `corntron`
The prebuilt `corntron` is an executable file, that composing the corns and runtimes,executing a command or application by corn config in composed environment. Released at [here](https://cnb.cool/viscropst/corntron/-/releases).
The artifact name of `corntron` command as `corntron_#{ARCH}.exe` for windows, `corntron_#{OS}_#{ARCH}` for others.
```
corntron --help
INF corntron version: <tagged version of prebuilt>
corntron.exe [options] <actions> [args]
actions was: [run-corn-config run-cmd run-corn]
options has:
  -cfg-base string
        /path/to/your/<corntron config folder> aka. base_dir of core.toml
  -corn-base string
        /path/to/your/<corns profiles folder>,for corn configs not in corntron conrig folder.
  -env-dirname string
        <folder name of env files to store>
  -mirror-type string
        mirror type, default is without mirror
  -no-waiting
        executing cryptron without waiting
  -rt-base string
        /path/to/your/<runtime profiles folder>,for runtrime configs not in corntron conrig folder.
  -running-base string
        /path/to/your/<corntron running folder>,for spliting the executables of corntron configs out of corntron config folder.
```
Usage of `run-cmd`
```
run-cmd [command] [args of command]
[command] is one of built-in comamnds and executable in `PATH` and absolute path of executable
[args of command] are args of command
if [command] and [args of command] are empty, will default to %COMPSPEC% or ${SHELL} or `/bin/sh`
```
Usage of `run-corn`
```
run-corn <corn config file name at corn config's env dirs,trimmed suffix `.toml`> [args of exec at corn config]
```
Usage of `run-corn-config`
```
INF corntron version: staging
run-corn-config [option] <corn config file suffixed `.toml.corn`>
Usage:
  -dir-as-base
        use the current corn file's dir as config base
error: flag: help requested
```

### How to build `corntron` command?
You can build corntron by using the `go build <source path>/cmd/corntron/main` command.
