#!/bin/sh
repowrok=$PWD/git
echo Start cloning bare of code
#ssh-keyscan e.coding.net >> ~/.ssh/known_hosts
#ssh -T git@e.coding.net
[ ! -d $repowrok ] && git clone --mirror https://cnb.cool/viscropst/corntron $repowrok
[ ! -d $repowrok ] && exit 1
cd $repowrok
echo Start push to GitCode
ssh-keyscan -t rsa gitcode.com >> ~/.ssh/known_hosts
ssh -T git@gitcode.com
git push --mirror git@gitcode.com:viscropst/corntron.git --force