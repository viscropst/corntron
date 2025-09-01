#!/bin/sh
repowrok=$PWD/git
echo Start cloning bare of code
ssh-keyscan e.coding.net >> ~/.ssh/known_hosts
ssh -T git@e.coding.net
[ ! -d $repowrok ] && git clone --mirror git@e.coding.net:visoft/imetnide/cryphtron-prototype.git $repowrok
[ ! -d $repowrok ] && exit 1
cd $repowrok
echo Start push to cnb.cool
git remote add target https://${CI_CNB_TOKEN}@cnb.cool/viscropst/corntron.git
git push --mirror target --force