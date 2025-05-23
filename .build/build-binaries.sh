#!/bin/sh
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
out_dir=$PWD/out
mkdir -p ${out_dir}
binaries="corntron"
os="windows 
linux 
freebsd 
openbsd 
darwin 
android 
ios"

cpu_arch="386
amd64
arm
arm64
loong64"
for sys in $os; do
    for arch in $cpu_arch; do
        platstr=${sys}"_"${arch}
        [ $sys == "ios" ] && continue
        [ $sys == "android" ] && continue
        [ $sys == "darwin" -a $arch != "amd64" -a $arch != "arm64" ] && continue
        [ $arch == "loong64" -a $sys != "linux" ] && continue
        export GOOS=$sys
        export GOARCH=$arch
        for bin in $binaries; do
            file_name=$bin"_"$arch
            [ $sys != "windows" ] && file_name=$file_name'_'$sys
            [ $sys == "windows" ] && file_name=$file_name'.exe'
            mkdir -p $out_dir
            echo "output file:" $out_dir/$bin/$file_name
            echo "build binary:" $bin
            CGO_ENABLED=0 go build -o $out_dir/$bin/$file_name $PWD/cmd/$bin/main || exit 1
        done
    done
done
