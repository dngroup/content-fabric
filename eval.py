#!/usr/bin/env python2
import argparse
import logging
import os
import sys
from argparse import RawTextHelpFormatter
from logging.handlers import RotatingFileHandler
from multiprocessing.pool import ThreadPool
from time import time, sleep

import jinja2
import numpy as np
import pandas as pd
import requests
from docker import Client

cli = Client(base_url='unix://var/run/docker.sock')
TARGET_DOCKER_COMPOSE_FILE = "docker-compose.cooked.yml"

cwd = os.path.basename(os.getcwd())
DOCKER_COMPOSE_NETWORK_NAME = "%s_default" % cwd.replace("-", "")
CHAIN_CODE_ID = "fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372"
PEER_COUNT = 5
CLIENT_COUNT = 10
ARRIVAL_TIME = 5
TE_COUNT = 10
CP_COUNT = 10
TE_PERCENT = 100
TE_PERCENT_PRICE = 100
CP_PERCENT = 100
CONSENSUS = "noops"
CONSENSUS_TIME_MAX = "100ms"  # 100ms
BATCH_SIZE=500

columns = ["peer_count",
           "client_count",
           "arrival_time",
           "te_count",
           "cp_count",
           "te_percent",
           "te_percent_price",
           "cp_percent",
           "consensus",
           "consensus_time_max",
           "batch_size",
           "max",
           "min",
           "mean",
           "res",
           "docker_server_diff",
           "date"]
filename = "data.csv"

parser = argparse.ArgumentParser(description='', epilog=
"""
""", formatter_class=RawTextHelpFormatter)
parser.add_argument('--chaincode_id', help='Chaincode id to use', default=CHAIN_CODE_ID)
parser.add_argument('--consensus', help='consensus to use default pbft', default=CONSENSUS)
parser.add_argument('--consensus_time_max', help='time max before consensus ', default=CONSENSUS_TIME_MAX)
parser.add_argument('--peer_count', type=int, help='Number of peer (need 3 or more)', default=PEER_COUNT)
parser.add_argument('--client_count', type=int, help='Number of coming client (need 1 or more)', default=CLIENT_COUNT)
parser.add_argument('--arrival_time', type=float, help='average time between 2 client', default=ARRIVAL_TIME)
parser.add_argument('--cp_count', type=int, help='Number of Content providers (need 1 or more)', default=CP_COUNT)
parser.add_argument('--cp_percent', type=int, help='percent of change to have content licensing in a CP',
                    default=CP_PERCENT)
parser.add_argument('--te_count', type=int, help='Number of Technical Enabler (need 1 or more)', default=TE_COUNT)
parser.add_argument('--te_percent', type=int, help='percent of change to have content in a TE', default=TE_PERCENT)
parser.add_argument('--te_percent_price', type=int,
                    help='percent of change to a better price thant the max price fixed by the CP',
                    default=TE_PERCENT_PRICE)
parser.add_argument('--batch_size', type=int, help='size of the batch for the pbft algo', default=500)
parser.add_argument('--no_run',action='store_false')

args = parser.parse_args()

CHAIN_CODE_ID = args.chaincode_id
PEER_COUNT = args.peer_count
CLIENT_COUNT = args.client_count
ARRIVAL_TIME = args.arrival_time
TE_COUNT = args.te_count
CP_COUNT = args.cp_count
TE_PERCENT = args.te_percent
TE_PERCENT_PRICE = args.te_percent_price
CP_PERCENT = args.cp_percent
CONSENSUS = args.consensus
CONSENSUS_TIME_MAX = args.consensus_time_max
BATCH_SIZE=args.batch_size
logger = logging.getLogger()

logger.setLevel(logging.DEBUG)

formatter = logging.Formatter('[%(asctime)s] p%(process)s {%(pathname)s:%(lineno)d} %(levelname)s - %(message)s',
                              '%m-%d %H:%M:%S')

file_handler = RotatingFileHandler('activity.log', 'a', 1000000, 1)

file_handler.setLevel(logging.DEBUG)
file_handler.setFormatter(formatter)
logger.addHandler(file_handler)

steam_handler = logging.StreamHandler()
steam_handler.setLevel(logging.DEBUG)
logger.addHandler(steam_handler)


def launch_user(networkName, vpnumber, chaincode):
    try:

        # docker run --net contentfabric_default -it -e USERID=lala -e CONTENTID=bbb.mp4 -e TIMEMAX=30  -e PEERADDR=172.20.0.14 -e EVENT_PORT=7053 -e REST_PORT=7050  -e CHAINCODE=fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372   -e USER=user_type1_0 dngroup/content-contract-user

        networking_config = cli.create_networking_config({
            networkName: cli.create_endpoint_config()
        })

        container = cli.create_container(image='dngroup/content-contract-user:reconect',
                                         environment={"PEERADDR": "vp%s" % (vpnumber),
                                                      # "EVENT_PORT": "%d" % (vpnumber + 3),
                                                      # "REST_PORT": "%d" % (vpnumber),
                                                      "CHAINCODE": chaincode,
                                                      "TIMEMAX": 30
                                                      }, networking_config=networking_config

                                         )
        start = time()
        cli.start(container=container.get("Id"))
        cli.wait(container=container.get("Id"))
        return int(cli.logs(container=container.get("Id")).split("price")[1][2:].split(",")[0]), time() - start
    except Exception as e:
        logging.error("failed to load time for container %s" % container.get("Id"))
        return -1, None


def waitChaincode(PEER_COUNT):
    number = 0
    t_end = time() + 60
    while number != PEER_COUNT and time() < t_end:
        containers = cli.containers(filters={"name": "dev-*"})
        number = len(containers)
        sleep(1)
    if time() > t_end:
        logging.error("Not all chainecode is deployed")
    return


def render(tpl_path, context):
    path, filename = os.path.split(tpl_path)
    return jinja2.Environment(
        loader=jinja2.FileSystemLoader(path or './')
    ).get_template(filename).render(context)


def register_chaincode(url="http://localhost:10000/chaincode"):
    payload = "{  \"jsonrpc\": \"2.0\",  \"method\": \"deploy\",  \"params\": {  \"type\": 1,  \"chaincodeID\": {  \"path\": \"https://github.com/dngroup/content-fabric/content-contract-cc\"  },  \"ctorMsg\": {  \"function\": \"init\",  \"args\": [  ]  }  },  \"secureContext\": \"admin\",    \"id\": 1 }"
    headers = {'content-type': "application/json"}

    response = requests.request("POST", url, data=payload, headers=headers)

    logging.debug(response.text)


context = {}
context = {"peer_count": PEER_COUNT,
           "te_count": TE_COUNT,
           "cp_count": CP_COUNT,
           "chaincode_id": CHAIN_CODE_ID,
           "te_percent": TE_PERCENT,
           "te_percent_price": TE_PERCENT_PRICE,
           "cp_percent": CP_PERCENT,
           "consensus": CONSENSUS,
           "consensus_time_max": CONSENSUS_TIME_MAX,
           "batch_size": BATCH_SIZE
           }

with open(TARGET_DOCKER_COMPOSE_FILE, "w") as f:
    f.write(render("./docker-compose.yml.tpl", context))


<<<<<<< HEAD
if not no_run:
=======
if not args.no_run:
>>>>>>> d45bbb7524a7ba6f60b588d0563cb5c05c760ecf
    try:
        os.system("docker-compose -f %s up -d " % TARGET_DOCKER_COMPOSE_FILE)
        logging.debug("waiting for the dockers to launch")
        sleep(3)
        logging.debug("deploying chaincode")
        register_chaincode()
        logging.debug("deploying chaincode [DONE]")
        # print("please press return when chaincode is everywhere")
        # raw_input()
        print("Waits for the chaincode to be compiled everywhere")
        waitChaincode(PEER_COUNT)
        print("Waits for the chaincode to be compiled everywhere [DONE]")

        gateway = \
            [item["IPAM"]["Config"][0]["Gateway"] for item in cli.networks() if
             item["Name"] == DOCKER_COMPOSE_NETWORK_NAME][0]

        closures = []

        rs = np.random.RandomState(1)


        def experiment(args):
            random_vp, timer = args
            sleep(timer)
            logging.debug("launching user %d" % timer)
            return timer, launch_user(DOCKER_COMPOSE_NETWORK_NAME, random_vp, CHAIN_CODE_ID)


        pool = ThreadPool(100)
        res = pool.map(experiment, zip(rs.randint(0, PEER_COUNT, CLIENT_COUNT),(np.cumsum(rs.poisson(ARRIVAL_TIME, CLIENT_COUNT)))))
        #map(experiment, zip(rs.randint(0, PEER_COUNT, CLIENT_COUNT), (np.cumsum(rs.poisson(ARRIVAL_TIME, CLIENT_COUNT)))))

        # get the number of chaincode server online
        containers = cli.containers(filters={"name": "dev-*"})
        diffonline = PEER_COUNT - len(containers)
        # save result
        try:
            # load the dataframe if it exists
            data = pd.DataFrame.from_csv(filename)
        except IOError as e:
            # otherwise, create it
            data = pd.DataFrame(columns=columns)

        resAsString = ', '.join(str(x) for x in res)
        # create a dataset containing the new data
        data_new = pd.DataFrame(np.array([[PEER_COUNT,
                                           CLIENT_COUNT,
                                           ARRIVAL_TIME,
                                           TE_COUNT,
                                           CP_COUNT,
                                           TE_PERCENT,
                                           TE_PERCENT_PRICE,
                                           CP_PERCENT,
                                           CONSENSUS,
                                           CONSENSUS_TIME_MAX,
                                           BATCH_SIZE,
                                           np.max([x[1][1] for x in res if x[1][1] is not None]),
                                           np.min([x[1][1] for x in res if x[1][1] is not None]),
                                           np.mean([x[1][1] for x in res if x[1][1] is not None]),
                                           resAsString,
                                           diffonline,
                                           time()
                                           ]]),
                                columns=columns)
        # add it to the old one, and save
        data = data.append(data_new)
        data.to_csv(filename)

        if str(res).find("None") > 0:
            sys.exit(-2)
            # raw_input()
    finally:
        print "stop docker"
        # raw_input()
        os.system("docker-compose -f %s kill -s 9" % TARGET_DOCKER_COMPOSE_FILE)
        os.system("docker-compose -f %s rm -f" % TARGET_DOCKER_COMPOSE_FILE)
