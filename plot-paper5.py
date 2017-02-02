#!/usr/bin/env python
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib import cm

fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

data = pd.DataFrame.from_csv("data.csv")

X, Y = np.meshgrid(sorted(list(set(data["te_count"].values))), sorted(list(set(data["client_count"].values))))
Z = np.empty((len(X), len(Y)))


def Zvalue(a, b, ax1="te_count", ax2="client_count", target="mean"):
    v = data[(data[ax1] == a) & (data[ax2] == b)][target].values
    if len(v) == 0:
        return np.nan
    else:
        return v[0]


Zvalue = np.vectorize(Zvalue)
Z = Zvalue(X, Y,target="max")
surf = ax.plot_wireframe(X, Y, Z, cmap=cm.jet, rstride=1, cstride=1)
ax.set_xlabel('TE/CP Count')
ax.set_ylabel('client count')
ax.set_zlabel('contract convergeance time')

plt.show()
