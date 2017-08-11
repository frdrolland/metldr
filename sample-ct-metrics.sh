#!/bin/bash

script=$(readlink -f $0)
scriptdir=$(dirname $script)

if [ "$#" -le 2 ]; then
  echo "Missing argument"
  echo "Usage: $0 <template file> <dest file>"
  exit 1
fi
sample_file=$1
dest_file=$2

echo "template\t:\t${sample_file}"
echo "template\t:\t${dest_file}"

current_ts=$(date +%s%N)



exit 0
