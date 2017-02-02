#!/usr/bin/env python

import argparse
import os
from argparse import RawTextHelpFormatter

import pandas as pd

CHAIN_CODE_ID = "fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372"

columns = ["peer_count",
           "client_count",
           "arrival_time",
           "te_count",
           "cp_count",
           "te_percent",
           "te_percent_price",
           "cp_percent",
           "consensus",
           "consensus_time_max"
           "do"
           ]
filename = "input.csv"

parser = argparse.ArgumentParser(description='', epilog=
"""
""", formatter_class=RawTextHelpFormatter)
parser.add_argument('--file', help='file with value to test', default=filename)

args = parser.parse_args()

# save result
try:
    # load the dataframe if it exists
    data = pd.DataFrame.from_csv(filename)
except IOError as e:
    # otherwise, create it
    data = pd.DataFrame(columns=columns)

for index, conf in data.iterrows():
    # print(conf["do"])
    if str(conf["do"]) == "True":
        # print("lala")
        continue
    # print (
    status = os.system(
        "./eval.py --peer_count %d --client_count %d --arrival_time %.2f --te_count %d --cp_count %d --te_percent %d --te_percent_price %d --cp_percent %d --consensus %s --consensus_time_max %s" % (
            conf["peer_count"], conf["client_count"], conf["arrival_time"], conf["te_count"], conf["cp_count"],
            conf["te_percent"],
            conf["te_percent_price"],
            conf["cp_percent"],
            conf["consensus"],
            conf["consensus_time_max"]))
    if os.WEXITSTATUS(status) == 0:
        data.iloc[index, data.columns.get_loc('do')] = "True"
        data.to_csv(filename)
    else:
        data.iloc[index, data.columns.get_loc('do')] = "Error"
        data.to_csv(filename)