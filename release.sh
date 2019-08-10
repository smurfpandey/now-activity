#!/bin/bash
# Script to build and create release package
# https://github.com/smurfpandey/now-activity

rm -rf now-activity
go build
tar -czvf now-activity.tar.gz now-activity now-activity.service now.smurfpandey.me.conf