#!/bin/bash

ssh -i ~/ssh_root_key root@$1 "systemctl stop wired && mount -o rw,remount /"

scp -i ~/ssh_root_key ./build/wired root@$1:/usr/bin/
#scp -i ~/ssh_root_key ./vector-gobot/build/libvector-gobot.so root@$1:/usr/lib/

scp -i ~/ssh_root_key -r ./webroot/* root@$1:/etc/wired/webroot/

#rsync -e "ssh -i ~/ssh_root_key" -avr ./modfiles/* root@$1:/etc/wired/mods/

ssh -i ~/ssh_root_key root@$1 "systemctl start wired"
