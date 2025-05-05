#!/bin/sh
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
out_dir=$PWD/out
mkdir -p ${out_dir}
binaries="cptron"
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
        echo "output file:" $out_dir/$file_name
        for bin in $binaries; do
            file_name=$bin"_"$arch
            [ $sys != "windows" ] && file_name=$file_name'_'$os
            [ $sys == "windows" ] && file_name=$file_name'.exe'
            mkdir -p $out_dir
            CGO_ENABLED=0 go build -o $out_dir/$file_name $PWD/cmd/$bin/main/main.go
        done
    done
done
