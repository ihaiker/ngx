#!/bin/bash

export BASE_PATH="$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )"

tag=$1

tag_commit=$(ls $BASE_PATH/../docs/versions/$tag.md)
if [ "$tag_commit" == "" ]; then
  echo "The tag $tag.md not found"
  exit 1
fi

tag_list=$(git tag|grep $tag)
if [ ! "$tag_list" == "" ]; then
    echo "Tag has created."
    exit 1
fi

cd $BASE_PATH/..

git tag $tag
git push --progress origin main --tags