#!/bin/bash

dest_file=$1

if [ -z "$dest_file" ]; then
  echo "Usage: $0 <log file>"
  exit 1
fi

#awk -F "|" '{ print gsub(/[\s\t]+$/, "", $6); }' ${dest_file} | sort -u
for stattype in $(awk -F "|" '{ print $6; }' ${dest_file} | tr -d " " | sort -u); do
  #echo "## ${stattype}"
  grep "${stattype}" ${dest_file} | awk -F "|" '{ print $NF;}'
done
