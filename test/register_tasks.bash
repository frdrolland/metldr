#!/bin/bash

script_dir="$(readlink -f $0)"
base_dir="$(dirname ${script_dir})"

export KAPACITOR_URL="http://localhost:9093"
database_name="testfro"
taskname="write_influxdb"


function check_retcode() {
  if [ $1 -ne 0 ]; then
      echo "[fatal] Command returned $1"
      exit $1
  fi
}

echo "${KAPACITOR_URL}"

echo "clean tasks before (re)create it"
kapacitor -url "${KAPACITOR_URL}" disable ${taskname}
kapacitor -url "${KAPACITOR_URL}" delete tasks ${taskname}

echo "Registering task for ${taskname} with script ${base_dir}/write_influxdb.tick"
kapacitor -url "${KAPACITOR_URL}" define write_influxdb \
    -type stream \
    -tick ${base_dir}/write_influxdb.tick \
    -dbrp ${database_name}.autogen
check_retcode $?

echo "Enabling task ${taskname}"
#kapacitor -url "${KAPACITOR_URL}" record stream -task ${taskname} -duration 20s
kapacitor -url "${KAPACITOR_URL}" enable ${taskname}
check_retcode $?

echo "Showing tasks"
#kapacitor -url "${KAPACITOR_URL}" record stream -task ${taskname} -duration 20s
kapacitor -url "${KAPACITOR_URL}" list tasks
check_retcode $?

echo "Changing log level to debug"
kapacitor -url "${KAPACITOR_URL}" level debug
check_retcode $?

kapacitor -url "${KAPACITOR_URL}" show ${taskname}

echo "[ok]"

