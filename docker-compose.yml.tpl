version: '2'
services:

{% for peer in range(0,peer_count): %}

  vp{{ peer}}:
    image: hyperledger/fabric-peer
    environment:
      - CORE_PEER_ID=vp{{ peer }}
      {% if peer != 0 %}
      - CORE_PEER_DISCOVERY_ROOTNODE=vp0:7051
      {% endif %}
      - CORE_PEER_ADDRESSAUTODETECT=true
      - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
      - CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft
      - CORE_PBFT_GENERAL_N={{ peer_count  }}
#      - CORE_LOGGING_LEVEL=DEBUG

    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: sh -c "peer node start "
    ports:
      - {{10000+peer*10}}-{{10000+peer*10+3}}:7050-7053
{% endfor %}


{% for te in range(0,te_count): %}
  te{{ te }}:
    image: dngroup/content-contract-te
    command: sh -c "sleep 5  &&  ./content-contract-te -events-address=vp{{ te % peer_count }}:7053 -events-from-chaincode={{ chaincode_id }} -TE-ID=te-{{te}}  -percent={{ te_percent }}  -percent-price={{ te_percent_price }} -rest-address=vp{{ te % peer_count }}:7050"
{% endfor %}



{% for cp in range(0,cp_count): %}
  cp{{ cp }}:
    image: dngroup/content-contract-cp
    command: sh -c "sleep 5 && ./content-contract-cp -events-address=vp{{ cp % peer_count }}:7053 -events-from-chaincode={{ chaincode_id }} -CP-ID=cp-{{cp}} -percent={{ cp_percent }} -rest-address=vp{{ cp % peer_count }}:7050"
{% endfor %}
