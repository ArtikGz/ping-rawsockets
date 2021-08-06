#!/bin/bash

# Remove build if there is an existing one
if [ -f ping ]; then
	sudo rm ping
fi

# This build and set suid permissions, so your are able to execute without sudo
go build -ldflags "-s -w" . && upx --ultra-brute ping && sudo chown root. ping && sudo chmod 4755 ping
echo -e "\nSuccessfully builded, now you can add it to path or move to /bin/ping (Care not override the default one)\nOr execute ./ping"
echo -e "Options:"
echo -e "\t-h [host] Host you want to ping (default: 127.0.0.1)"
echo -e "\t-c [times] Number of packets you want to send (default: 4)"

