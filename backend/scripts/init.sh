#!/bin/bash

# Cron jobs
service cron start
# Go rebuilder
CompileDaemon -build="go build -o /build/go/src/backend" -command="/build/go/src/backend"
