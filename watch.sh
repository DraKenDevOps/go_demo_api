#!/bin/bash
while true; do
  # Run the app in the background
  go run main.go &
  PID=$!
  # Wait for a file change in the current directory
  inotifywait -e modify -r .
  # Kill the previous process and restart
  kill $PID
done
