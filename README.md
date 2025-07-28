# Project corntron
Running an extended TOML-defined environment with one or more extended TOML-defined configurations.  
This source code is a reference implementation and prototype design of corntron. 
It is licensed under Mulan PSL v2.

## What is a corn?
A corn is a program that can be run with a set of arguments and a set of environment variables.
- It's defined by one or more extended TOML files.
- Can be defined to run a program.
- Can be defined as a meta-configuration for a corn.

## What is a runtime?
A runtime is a set of environment variables and a set of configuration variables.
- It's defined by an extended TOML file.
- Can be defined to prepare an environment for corn or interactive commandline.
- Can be defined to prepare public tools (eg. jdk, mvn, python) for a corn or environment for interactive commandline.
- Can be used to automatically do some work when you need mirror settings because of country-specific networking issues.

## What is directory structure?
- `corns` directory contains corn environments and the corn environment root for common platform.
- `corns/_envs/<corn-name>.toml` contains corn configurations.
- `runtimes` directory contains runtimes and the runtime environment root for common platform.
- `runtimes/_envs/<runtime-name>.toml` directory contains runtimes config.
- `#{GOOS}-#{GOARCH}` directory contains current operating system and architecture specific binaries of corns and runtimes. 
- `#{GOOS}-#{GOARCH}/corns` directory contains current operating system and architecture specific corns.
- `#{GOOS}-#{GOARCH}/runtimes` directory contains current operating system and architecture specific runtimes.
> [`GOOS`](https://github.com/golang/go/tree/master/src/internal/goos) and [`GOARCH`](https://github.com/golang/go/tree/master/src/internal/goarch) are defined by Go programming Language. These will be replaced by LLVM Triple style [`os`](https://github.com/llvm/llvm-project/blob/23d4756c4bfce06a98c9c03b24752d32760ac22b/llvm/include/llvm/TargetParser/Triple.h#L46) and [`arch`](https://github.com/llvm/llvm-project/blob/23d4756c4bfce06a98c9c03b24752d32760ac22b/llvm/include/llvm/TargetParser/Triple.h#L46) in the future.

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
The corn config file is a corntron environment config file that is used to running a program or command. 
```toml
# This is a corn config file.

# `meta_only` defines whether this corn config is a meta-corn config.
# meta-corn config is a corn config that not executing the commands in `bootstrap_exec`.
# If `meta_only` is true, the corn config is a meta-corn config.
# Otherwise, the corn config is a normal corn config.
meta_only = false

# `depend_corns` is the set of corns that this corn depends on.
# If the depend corns are not running, this corn will not be executed.
# And the `vars` and `envs` of this corn will be merged with the `vars` and `envs` of the depend corns.
# The value of `depend_corns` is a list of string.
depend_corns = ["foo"]

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
The runtime environment config file is used to define a runtime environment for a corntron program (a corn config file or interactive commandline).
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

## How to use corntron?
You can integrate corntron into your program by using the `corntron` package.
Or you can use the prebuilt [`corntron` command](https://cnb.cool/viscropst/corntron/-/releases) to run a corntron program.

## How to build corntron command?
You can build corntron by using the `go build cmd/corntron/main.go` command.
Or you can use the `go install github.com/viscropst/corntron/cmd/corntron` command.
