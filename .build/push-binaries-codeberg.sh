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
    "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/tags/"${version}`
echo "After get release by tag"
echo "result:" $result
release_id=`echo $result | jq -r '.id'`
if [ "${release_id}" != "" ]; then
    echo "Release already exists, deleting"
    curl -s -X DELETE \
        -H "accept: application/json" \
        -H "authorization: token "${CI_CB_TOKEN} \
        "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/"${release_id}
fi
echo "Creating the CodeBerg release, tagged:"${version}
release_id=`curl -s -X POST \
    -H 'accept: application/json' \
    -H "Content-Type: application/json" \
    -H "authorization: token "${CI_CB_TOKEN} \
    -d '{"tag_name":"'${version}'","name":"'${version}' release","body":"this is a release by tag","draft":false,"prerelease":'${pre_release}'}' \
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
            curl -X POST \
                -F "attachment=@${file}" \
                -H "authorization: token "${CI_CB_TOKEN} \
                "https://"${GENERIC_ARTIFACTS_MACHINE}${CODEBERG_RELEASE_PATH}"/"${release_id}"/assets"
        fi
        echo "Deleting from generic package store"
        curl -s -X DELETE \
            -H "authorization: token "${CI_CB_TOKEN} \
            "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename} || true
        echo "Uploading to generic package store"
        echo "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename}
        curl -T $file \
            -H "authorization: token"${CI_CB_TOKEN} \
            "https://"${GENERIC_ARTIFACTS_MACHINE}${GENERIC_ARTIFACTS_STORE}"/"${binary}"/"${version}"/"${filename}
    done
done