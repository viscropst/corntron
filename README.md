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
- `corns` directory contains corns.
- `corns/_envs/<corn-name>.toml` directory contains corns config.
- `runtimes` directory contains runtimes.
- `runtimes/_envs/<runtime-name>.toml` directory contains runtimes config.
- `#{GOOS}-#{GOARCH}` directory contains current operating system and architecture specific binaries of corns and runtimes. 
- `#{GOOS}-#{GOARCH}/corns` directory contains current operating system and architecture specific corns.
- `#{GOOS}-#{GOARCH}/runtimes` directory contains current operating system and architecture specific runtimes.
> [`GOOS`](https://github.com/golang/go/tree/master/src/internal/goos) and [`GOARCH`](https://github.com/golang/go/tree/master/src/internal/goarch) are defined by Go programming Language. These will be replaced by LLVM Triple style [`os`](https://github.com/llvm/llvm-project/blob/23d4756c4bfce06a98c9c03b24752d32760ac22b/llvm/include/llvm/TargetParser/Triple.h#L46) and [`arch`](https://github.com/llvm/llvm-project/blob/23d4756c4bfce06a98c9c03b24752d32760ac22b/llvm/include/llvm/TargetParser/Triple.h#L46) in the future.

## How to use corntron?
You can integrate corntron into your program by using the `corntron` package.
Or you can use the prebuilt `corntron` command to run a corntron program.

## How to build corntron command?
You can build corntron by using the `go build cmd/corntron/main.go` command.
Or you can use the `go install github.com/viscropst/corntron/cmd/corntron` command.
