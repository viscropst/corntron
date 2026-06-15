#!/bin/sh

sh .build/mirror-preparing.sh
sh .build/mirror-to-github.sh

out_dir=$PWD/out
version="staging"
pre_release="true"
[ ${CNB_BRANCH} != "main" ] && version=${CNB_BRANCH}
[ ${CNB_BRANCH} != "main" ] && pre_release="false"

apk add --no-cache curl jq

GENERIC_ARTIFACTS_MACHINE="api.github.com"
GITHUB_RELEASE_PATH="/repos/viscropst/corntron/releases"
GENERIC_UPLOADS_MACHINE="uploads.github.com"

authorization_token="Bearer "${CI_GH_TOKEN}
application_type="application/vnd.github+json"

release_id=""
echo "Getting release by tag"
result=`curl -s -L -H "authorization: ${authorization_token}" \
    -H "accept: ${application_type}" \
    -H "X-GitHub-Api-Version: 2026-03-10" \
    "https://"${GENERIC_ARTIFACTS_MACHINE}${GITHUB_RELEASE_PATH}"/tags/"${version}`
echo "After get release by tag"
echo "result:" $result
release_id=`echo $result | jq -r '.id'`
if [ "${release_id}" != "null" ]; then
    echo "Release already exists, deleting"
    curl -s -L -X DELETE \
        -H "accept: "${application_type} \
        -H "authorization: ${authorization_token}" \
        -H "X-GitHub-Api-Version: 2026-03-10" \
        "https://"${GENERIC_ARTIFACTS_MACHINE}${GITHUB_RELEASE_PATH}"/"${release_id}
fi
echo "Creating the Github release, tagged:"${version}

[ ! -d "./out/cortron" ] && mkdir -p ./out/cortron

echo ${LATEST_CHANGELOG} > ./out/cortron/CHANGELOG.md

release_id=`curl -s -L -H "Accept: ${application_type}" \
    -H "Authorization: ${authorization_token}" \
    -H "X-GitHub-Api-Version: 2026-03-10"  \
    "https://"${GENERIC_ARTIFACTS_MACHINE}${GITHUB_RELEASE_PATH} \
    -X POST \
    -d '{"tag_name":"'${version}'","name":"'${version}'","body":"","draft":false,"generate_release_notes":true,"prerelease":'${pre_release}'}' | jq -r ".id"`

if [ "${release_id}" == "null" ];then 
    echo "Failed to get or create release, exit"
    exit 1
fi

for binary in `ls ${out_dir}/`;
do
    for filename in `ls ${out_dir}/${binary}`;
    do
        artifact_file=${binary}"/"${filename}
        file=${out_dir}"/"${artifact_file} 
        echo "now pushing" $file "with version:" $version
        echo "artifact file was:" ${artifact_file}
        if [ "${release_id}" != "null" ]; then
            echo "Uploading to Github release"
            echo "https://"${GENERIC_UPLOADS_MACHINE}${GITHUB_RELEASE_PATH}"/"${release_id}"/assets?name="${filename}
            curl -L --max-time 120 -X POST \
                -H "authorization: ${authorization_token}" \
                -H "X-GitHub-Api-Version: 2026-03-10" \
                -H "Content-Type: application/octet-stream" \
                "https://"${GENERIC_UPLOADS_MACHINE}${GITHUB_RELEASE_PATH}"/"${release_id}"/assets?name="${filename} \
                --data-binary "@"${file}
        fi
    done
done