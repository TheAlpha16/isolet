#!/bin/sh

adduser -D level1

echo "level1:$USER_PASSWORD" | chpasswd
echo "hi this is level 1: $FLAG" > /home/level1/flag.txt
/usr/sbin/sshd -D