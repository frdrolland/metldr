#!/bin/bash

dest_file=$1

if [ -z "$dest_file" ]; then
  echo "Usage: $0 <log file>"
  exit 1
fi

for stattype in $(awk -F "|" '{ print $6; }' ${dest_file} | tr -d " " | sort -u); do
  grep "${stattype}" ${dest_file} | awk -F "|" '{ print $NF;}'
done
