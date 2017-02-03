#!/usr/bin/env python
# !/usr/bin/env python

import argparse
from argparse import RawTextHelpFormatter

import numpy as np

CHAIN_CODE_ID = "fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372"
FILENAME = "input.csv"
MAX_USER = 600
STEPS_USERS = 5
MIN_USER = 6
MIN_PEERS = 4
MAX_PEERS = 10
STEPS_PEERS = 1
MAX_TE = 600
STEPS_TES = 5
MIN_TE = 1
CONSENSUS_TIME_MAX = "100ms"

PEER_COUNT = 4

TE_PERCENT = 100
TE_PERCENT_PRICE = 100
CP_PERCENT = 100
CONSENSUS = "noops"

CP_COUNT = 10

parser = argparse.ArgumentParser(description='', epilog=
"""
""", formatter_class=RawTextHelpFormatter)
parser.add_argument('--file', help='file to save', default=FILENAME)
parser.add_argument('--peer_count', help='peer_count ', default=PEER_COUNT)
parser.add_argument('--max_user', type=int, help='max user in 1 min', default=MAX_USER)
parser.add_argument('--min_user', type=int, help='min user in 1 min', default=MIN_USER)
parser.add_argument('--steps_users', type=int, default=STEPS_USERS)
parser.add_argument('--steps_tes', type=int, default=STEPS_TES)
parser.add_argument('--max_te', type=int, help='max number of  cp and te', default=MAX_TE)
parser.add_argument('--min_te', type=int, help='min number of cp and te', default=MIN_TE)

parser.add_argument('--max_peers', type=int, help='max number of  cp and te', default=MAX_PEERS)
parser.add_argument('--min_peers', type=int, help='min number of cp and te', default=MIN_PEERS)
parser.add_argument('--step_peers', type=int, help='min number of cp and te', default=STEPS_PEERS)
parser.add_argument('--consensus_time_max', help='time max before consensus ', default=CONSENSUS_TIME_MAX)
parser.add_argument('--te_percent', type=int, help='percent of change to have content in a TE', default=TE_PERCENT)
parser.add_argument('--te_percent_price', type=int,
                    help='percent of change to a better price thant the max price fixed by the CP',
                    default=TE_PERCENT_PRICE)
parser.add_argument('--cp_percent', type=int, help='percent of change to have content licensing in a CP',
                    default=CP_PERCENT)
parser.add_argument('--cp_count', type=int, help='Number of Content providers (need 1 or more)', default=CP_COUNT)
parser.add_argument('--consensus', help='consensus to use default pbft', default=CONSENSUS)

args = parser.parse_args()
MAX_USER = args.max_user
STEPS_USERS = args.steps_users
STEPS_TES = args.steps_tes
MIN_USER = args.min_user

MAX_PEERS = args.max_peers
MIN_PEERS = args.min_peers
STEPS_PEERS = args.step_peers

MAX_TE = args.max_te
MIN_TE = args.min_te
PEER_COUNT = args.peer_count
CONSENSUS_TIME_MAX = args.consensus_time_max
TE_PERCENT = args.te_percent
TE_PERCENT_PRICE = args.te_percent_price
CP_PERCENT = args.cp_percent
CONSENSUS = args.consensus

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
           "do"
           ]

users = np.geomspace(MIN_USER, MAX_USER, STEPS_USERS)
cptes = np.geomspace(MIN_TE, MAX_TE, STEPS_TES)
peers = list(sorted(set([int(x) for x in np.linspace(MIN_PEERS,MAX_PEERS,STEPS_PEERS)])))


print "," + ",".join(columns)
for user in users:
    for te in cptes:
        for peer in peers:
            print ("0,%d ,%d , %.2f, %d ,%d , %d , %d , %d , %s , %s, False" % (
                peer,
                user,
                60.0 / user,
                te,
                CP_COUNT,
                TE_PERCENT,
                TE_PERCENT_PRICE,
                CP_PERCENT,
                CONSENSUS,
                CONSENSUS_TIME_MAX
            ))
