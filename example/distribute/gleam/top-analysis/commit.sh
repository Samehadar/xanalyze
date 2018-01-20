#!/usr/bin/env bash

echo 'commit message：'

read msg

git pull origin master

git add .

git commit -m '$msg'

git push  origin master

