#!/bin/bash


delete_part() {
  if [[ "$1" == "" ]] ; then
    echo "delete_part need parameter"
    exit -1
  fi
  s=$1
  find . -type f -name "*.go" -exec perl -i -p0e "s/^[\t ]*\/\/\s*begin\s+${s}\b.*$\s*[\s\S]*?\/\/\s*end\s+${s}\b.*$//gm" {} \;
}

delete_margin() {
  if [[ "$1" == "" ]] ; then
    echo "delete_margin need parameter"
    exit -1
  fi
  s=$1
  find . -type f -name "*.go" -exec perl -i -pe "s/^[\t ]*\/\/\s*(begin|end)\s+${s}\b.*$//g" {} \;
}

usage() {
  echo "Usage:"
  echo "  $0 <new-name> [addin]"
  echo ""
  echo "Available addin:"
  echo "  basic  [default]"
  echo "  gin"
  echo "  tcp"
}

if [[ "$1" == "" ]]; then
  usage
  exit -1
fi

new_name=$1
if [[ "$new_name" == "basic" || "$new_name" == "gin" || "$new_name" == "tcp" ]]; then
  echo "new-name cannot be basic/gin/tcp, maybe you can use test_${new_name}"
  echo ""
  usage
  exit -1
fi

addin=$2
if [[ "$addin" == "" ]]; then
  addin="basic"
fi

if [[ "$addin" == "gin" ]]; then
  delete_part basic
  delete_part tcp
  delete_margin gin
elif [[ "$addin" == "tcp" ]]; then
  delete_part basic
  delete_part gin
  delete_margin tcp
elif [[ "$addin" == "basic" ]]; then
  delete_part gin
  delete_part tcp
  delete_margin basic
else
  echo "Invalid addin: $addin"
  echo ""
  usage
  exit -1
fi

find . -type f \( -name "*.go" -or -name "go.mod" -or -name "*.yaml" \) -exec sed -i "s/TEMPLATE/${new_name}/g" {} \;
