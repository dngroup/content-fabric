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
      - CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN={{ consensus  }}
      - CORE_PBFT_GENERAL_N={{ peer_count  }}
      - CORE_PBFT_GENERAL_BATCHSIZE={{ batch_size }}
      - CORE_NOOPS_BLOCK_WAIT={{ consensus_time_max }}
      - CORE_PBFT_GENERAL_TIMEOUT_BATCH={{ consensus_time_max }}
#      - CORE_LOGGING_LEVEL=DEBUG
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    {% if peer == 0 %}
    command: sh -c "peer node start -v"
    {% else %}
    depends_on:
        - vp0
    command: sh -c "sleep 2 && peer node start -v"
    {% endif %}
    ports:
      - {{10000+peer*10}}-{{10000+peer*10+3}}:7050-7053
{% endfor %}


{% for te in range(0,te_count): %}
  te{{ te }}:
    image: dngroup/content-contract-te
    depends_on:
    {% for peer in range(0,peer_count): %}
        - vp{{ peer }}
    {% endfor %}
    # start te one by one
    #{% if te != 0 %}
    #    - te{{ te-1 }}
    #{% endif %}
    command: sh -c "sleep 5  &&  ./content-contract-te -events-address=vp{{ te % peer_count }}:7053 -events-from-chaincode={{ chaincode_id }} -TE-ID=te-{{te}}  -percent={{ te_percent }}  -percent-price={{ te_percent_price }} -rest-address=vp{{ te % peer_count }}:7050"
{% endfor %}



{% for cp in range(0,cp_count): %}
  cp{{ cp }}:
    image: dngroup/content-contract-cp
    depends_on:
    {% for peer in range(0,peer_count): %}
        - vp{{ peer }}
    {% endfor %}
    command: sh -c "sleep 5 && ./content-contract-cp -events-address=vp{{ cp % peer_count }}:7053 -events-from-chaincode={{ chaincode_id }} -CP-ID=cp-{{cp}} -percent={{ cp_percent }} -rest-address=vp{{ cp % peer_count }}:7050"
{% endfor %}
