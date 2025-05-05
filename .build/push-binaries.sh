#!/bin/sh

out_dir=$PWD/out
version="staging"
[ ${CODING_BRANCH} -neq "main" ] && version=${CODING_BRANCH}

GENERIC_ARTIFACTS_MACHINE="visoft-generic.pkg.coding.net"
GENERIC_ARTIFACTS_STORE="/imetnide/binary"
GENERIC_ARTIFACTS_USER=binary-1746428262480
GENERIC_ARTIFACTS_PWD=9dfcfc87722fdc68c7883d2f82c314ffb6135a93

for binary in `ls ${out_dir}/`;
do
    for filename in `ls ${out_dir}/${binary}`;
    do
        file=${out_dir}"/"${filename}
        artifact_file=${binary}"/"${filename}
        echo "now pushing" $file "with version:" $version
        echo "artifact file was:" ${artifact_file}
        curl -T $file \
            -u ${GENERIC_ARTIFACTS_USER}:${GENERIC_ARTIFACTS_PWD} \
            "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${artifact_file}"?version="${version}
    done
done