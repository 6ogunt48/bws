#!/bin/bash

URL="http://localhost/"
REQUEST_COUNT=100000
for ((i=1; i<=REQUEST_COUNT; i++))
do
  curl -s $URL &
  sleep 0.01
done

wait
echo "Done sending $REQUEST_COUNT requests"
