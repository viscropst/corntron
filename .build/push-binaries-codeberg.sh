#!/bin/sh

sh .build/mirror-preparing.sh
sh .build/mirror-to-codeberg.sh

out_dir=$PWD/out
version="staging"
pre_release="true"
[ ${CNB_BRANCH} != "main" ] && version=${CNB_BRANCH}
[ ${CNB_BRANCH} != "main" ] && pre_release="false"

apk add --no-cache curl jq

GENERIC_ARTIFACTS_MACHINE="codeberg.org/api/"
CODEBERG_RELEASE_PATH="v1/repos/viscropst/corntron/releases"
GENERIC_ARTIFACTS_STORE="packages/viscropst/generic"


release_id=""
echo "Getting release by tag"
result=`curl -s -H "authorization: token "${CI_CB_TOKEN} \
    -H 'accept: application/json' \
    "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/tags/"${version}`
echo "After get release by tag"
echo "result:" $result
release_id=`echo $result | jq -r '.id'`
if [ "${release_id}" != "null" ]; then
    echo "Release already exists, deleting"
    curl -s -X DELETE \
        -H "accept: application/json" \
        -H "authorization: token "${CI_CB_TOKEN} \
        "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/"${release_id}
fi
echo "Creating the CodeBerg release, tagged:"${version}

echo ${LATEST_CHANGELOG} > out/corntron/CHANGELOG.md
release_id=`curl -s -X POST \
    -H 'accept: application/json' \
    -H "Content-Type: application/json" \
    -H "authorization: token "${CI_CB_TOKEN} \
    -d '{"tag_name":"'${version}'","name":"'${version}' release","body":"read changelog in CHANGELOG.md","draft":false,"prerelease":'${pre_release}'}' \
    "https://${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}" | jq -r ".id"`

for binary in `ls ${out_dir}/`;
do
    for filename in `ls ${out_dir}/${binary}`;
    do
        artifact_file=${binary}"/"${filename}
        file=${out_dir}"/"${artifact_file} 
        echo "now pushing" $file "with version:" $version
        echo "artifact file was:" ${artifact_file}
        if [ "${release_id}" != "null" ]; then
            echo "Uploading to CodeBerg release"
            echo "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/"${release_id}"/assets"
            curl --max-time 120 -X POST \
                -H "authorization: token "${CI_CB_TOKEN} \
                "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/"${release_id}"/assets" \
                -F "attachment=@${file}"
        fi
        echo "Deleting from generic package store"
        curl -s --max-time 120 -X DELETE \
            -H "authorization: token "${CI_CB_TOKEN} \
            "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename}
        echo "\nUploading to generic package store"
        echo "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename}
        curl -H "authorization: token "${CI_CB_TOKEN} \
            "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename} \
            -T $file --max-time 120
    done
done