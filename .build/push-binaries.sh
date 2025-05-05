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
    curl -T ${binary} \
        -u ${GENERIC_ARTIFACTS_USER}:${GENERIC_ARTIFACTS_PWD} \
        "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"?version="${version}
done