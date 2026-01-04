#!/bin/sh
repowrok=$PWD/git
echo Start cloning bare of code
#ssh-keyscan e.coding.net >> ~/.ssh/known_hosts
#ssh -T git@e.coding.net
[ ! -d $repowrok ] && git clone --mirror https://cnb.cool/viscropst/corntron.git $repowrok
[ ! -d $repowrok ] && exit 1
cd $repowrok
echo Start push to CodeBerg
echo "$(nslookup -type=a codeberg.org 183.60.83.19 | grep Address | tail -n 1 | cut -b 10- ) codeberg.org" >> /etc/hosts
ssh-keyscan -t rsa codeberg.org >> ~/.ssh/known_hosts
ssh -T git@codeberg.org
git push --mirror ssh://git@codeberg.org/viscropst/corntron.git --force