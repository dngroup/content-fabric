#!/usr/bin/env python
import argparse

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib import cm
from mpl_toolkits.mplot3d import Axes3D

dir(Axes3D)

parser = argparse.ArgumentParser()

parser.add_argument('--axis1', default="te_count")
parser.add_argument('--axis2', default="client_count")
parser.add_argument('--axis3', default="min")

args = parser.parse_args()

fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

data = pd.DataFrame.from_csv("data.csv")

assert args.axis1 in list(data)
assert args.axis2 in list(data)
assert args.axis3 in list(data)

X, Y = np.meshgrid(sorted(list(set(data[args.axis1].values))), sorted(list(set(data[args.axis2].values))))
Z = np.empty((len(X), len(Y)))


def Zvalue(a, b, ax1, ax2, ax3):
    v = data[(data[ax1] == a) & (data[ax2] == b)][ax3].values
    if len(v) == 0:
        return np.nan
    else:
        return v[0]


Zvalue = np.vectorize(Zvalue)
Z = Zvalue(X, Y, args.axis1, args.axis2, args.axis3)
surf = ax.plot_wireframe(X, Y, Z, cmap=cm.jet, rstride=1, cstride=1)
ax.set_xlabel(args.axis1)
ax.set_ylabel(args.axis2)
ax.set_zlabel(args.axis3)

plt.show()
