#!/usr/bin/env bash
On_White='\033[47m'
RED='\033[0;31m'
NC='\033[0m' # No Color
MASTERIP=$(dig +short myip.opendns.com @resolver1.opendns.com)
echo -e "${RED}${On_White}master IP: ${MASTERIP}${NC}"

sed -i.bak '/jessie-backports/d' /etc/apt/sources.list
sed -i.bak '/stretch/d' /etc/apt/sources.list
echo 'deb http://httpredir.debian.org/debian jessie-backports main' >> /etc/apt/sources.list
echo 'deb http://httpredir.debian.org/debian stretch main' >> /etc/apt/sources.list
apt-get update -q
DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -qy salt-master salt-common
echo "fileserver_backend:" >> /etc/salt/master
echo "  - roots" >> /etc/salt/master
echo "  - minion" >> /etc/salt/master
echo "  - git" >> /etc/salt/master
echo "" >> /etc/salt/master
echo "file_recv: True" >> /etc/salt/master
echo "auto_accept: True" >> /etc/salt/master
service salt-master restart
