#!/bin/sh
repowrok=$PWD/git
id_name=id_rsa
[ ! -f /usr/bin/git ] && sed -i 's/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g' /etc/apk/repositories
[ ! -f /usr/bin/git ] && apk add --no-cache --update git gpg less openssh patch perl base64
[ ! -d ~/.ssh ] && mkdir ~/.ssh
echo ${CI_PUBLIC_KEY} > ~/.ssh/${id_name}.pub
echo public key was:
cat ~/.ssh/${id_name}.pub

echo ${CI_PRIVATE_KEY} | base64 -d > ~/.ssh/${id_name}
echo private key was:
# cat ~/.ssh/${id_name}
# echo \n
chmod 600 ~/.ssh/${id_name}

[ ! -f ~/.ssh/known_hosts ] && touch ~/.ssh/known_hosts
git config --global user.email "viscropst@petalmail.com"
git config --global user.name "viscropst"