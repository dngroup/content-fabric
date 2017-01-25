#!/usr/bin/env bash
sed -i.bak '/jessie-backports/d' /etc/apt/sources.list
sed -i.bak '/stretch/d' /etc/apt/sources.list
echo 'deb http://httpredir.debian.org/debian jessie-backports main' >> /etc/apt/sources.list
echo 'deb http://httpredir.debian.org/debian stretch main' >> /etc/apt/sources.list
apt-get update -q
DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -qy salt-minion salt-common
wget -qO- https://get.docker.com/ | sh

echo $2 > /etc/salt/minion_id

if [ $# -eq 0 ]
  then
    echo "you sould to do that :"
    echo '$MASTER salt >> /etc/hosts'
else
    echo $1 salt >> /etc/hosts
fi
