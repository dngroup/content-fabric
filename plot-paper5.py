#!/usr/bin/env python
import argparse
from scipy import interpolate
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib import cm
from mpl_toolkits.mplot3d import Axes3D

dir(Axes3D)

parser = argparse.ArgumentParser()

parser.add_argument('--axis1', default="consensus_time_max")
parser.add_argument('--axis2', default="batch_size")
parser.add_argument('--axis3', default="mean")
parser.add_argument('--resolution', default=80,type=int)

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
        return np.mean(v)


Zvalue = np.vectorize(Zvalue)
Z = Zvalue(X, Y, args.axis1, args.axis2, args.axis3)


#smoothing


#xnew, ynew = np.mgrid[1:39:80j, 1:79:80j]
xnew, ynew = np.mgrid[np.min(X):np.max(X):args.resolution*1j, np.min(Y):np.max(Y):args.resolution*1j]
tck = interpolate.bisplrep(X, Y, Z, s=0)
znew = interpolate.bisplev(xnew[:,0], ynew[0,:], tck)

ax = fig.gca(projection='3d')
ax.plot_surface(xnew, ynew, znew, cmap='summer', rstride=1, cstride=1, alpha=None, antialiased=True)



#ax = fig.gca(projection='3d')
#ax.plot_surface(xnew, ynew, znew, cmap='summer', rstride=1, cstride=1, alpha=None, antialiased=True)


#surf = ax.plot_wireframe(X, Y, Z,  cmap='summer', rstride=1, cstride=1, alpha=None, antialiased=True)
#ax.plot_surface(xnew, ynew, znew, cmap='summer', rstride=1, cstride=1, alpha=None, antialiased=True)
ax.set_xlabel(args.axis1)
ax.set_ylabel(args.axis2)
ax.set_zlabel(args.axis3)

plt.show()
#plt.savefig('foo%03d.png'%args.resolution, bbox_inches='tight')
