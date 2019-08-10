#!/bin/bash
# Script to build and create release package
# https://github.com/smurfpandey/whats-playing

rm -rf whats-playing
go build
tar -czvf whats-playing.tar.gz whats-playing whats-playing.service now.music.smurfpandey.me.conf