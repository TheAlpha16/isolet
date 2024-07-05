#!/bin/sh

adduser -D level2

echo "level2:$USER_PASSWORD" | chpasswd
echo "hi this is level 2: $FLAG" > /etc/flag.txt
chmod +r /etc/flag.txt
/usr/sbin/sshd -D