#!/usr/bin/env bash
# docker swarm init --advertise-addr eth1
# salt '*' cmd.run "docker swarm join --token SWMTKN-1-22pt6tpduvk7bkgtwn2xi0mkazd6uhox9hhzl9e1n6iriddagy-9w5nwencev2nww6nnurelesya 192.168.56.101:2377"

salt 'node1' cmd.run "docker run -d \
      --name peer1
      -e CORE_PEER_ID=peer1 \
      -e CORE_PEER_ADDRESSAUTODETECT=true \
      -e CORE_VM_ENDPOINT=unix:///var/run/docker.sock \
      -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 7050-7053:7050-7053 \
      dngroup/fabric-peer sh -c \"peer node start\""

salt 'node2' cmd.run "docker run -d \
      --name peer2
      -e CORE_PEER_ID=peer2 \
      -e CORE_PEER_DISCOVERY_ROOTNODE=${1}:7051
      -e CORE_PEER_ADDRESSAUTODETECT=true \
      -e CORE_VM_ENDPOINT=unix:///var/run/docker.sock \
      -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 7050-7053:7050-7053 \
      dngroup/fabric-peer sh -c \"peer node start\""

salt 'node3' cmd.run "docker run -d \
      --name peer3
      -e CORE_PEER_ID=peer3 \
      -e CORE_PEER_DISCOVERY_ROOTNODE=${1}:7051
      -e CORE_PEER_ADDRESSAUTODETECT=true \
      -e CORE_VM_ENDPOINT=unix:///var/run/docker.sock \
      -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 7050-7053:7050-7053 \
      dngroup/fabric-peer sh -c \"peer node start\""

salt 'node4' cmd.run "docker run -d \
      --name peer4
      -e CORE_PEER_ID=peer4 \
      -e CORE_PEER_DISCOVERY_ROOTNODE=${1}:7051
      -e CORE_PEER_ADDRESSAUTODETECT=true \
      -e CORE_VM_ENDPOINT=unix:///var/run/docker.sock \
      -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 7050-7053:7050-7053 \
      dngroup/fabric-peer sh -c \"peer node start\""

salt 'node5' cmd.run "docker run -d \
      --name peer5
      -e CORE_PEER_ID=peer5 \
      -e CORE_PEER_DISCOVERY_ROOTNODE=${1}:7051
      -e CORE_PEER_ADDRESSAUTODETECT=true \
      -e CORE_VM_ENDPOINT=unix:///var/run/docker.sock \
      -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 7050-7053:7050-7053 \
      dngroup/fabric-peer sh -c \"peer node start\""

      salt '*' cmd.script_retcode salt://scripts/runme.sh 'arg1 arg2 "arg 3"'
