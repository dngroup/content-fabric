import os
import time
import numpy as np



repeat=1
client_count=10
arrival_time=1
#batch_candidates=np.round(np.logspace(1,3,10))
#consensus_time_candidates=np.round(np.logspace(0,2.5,10))
batch_candidates=np.round(np.linspace(1,500,6))
consensus_time_candidates=np.round(np.linspace(1,1000,6))
vp_delay=10

if os.path.exists("data.csv"):
  os.rename("data.csv", "data_%d.csv"%time.time())

d="sudo service docker restart\nsudo ./eval.py --peer_count 4 --client_count %d --arrival_time %d --te_count %d --cp_count %d --te_percent 100 --te_percent_price 100 --cp_percent 100 --consensus  pbft  --consensus_time_max  %dms --batch_size %d --vp_delay=%d"

for i in range(0,repeat):
  print("echo 'RUN %d'" %i)
  for consensus_time_max in consensus_time_candidates:
    for batch_size in batch_candidates:
      print(d%(client_count,arrival_time,5,5,consensus_time_max,batch_size,vp_delay))
