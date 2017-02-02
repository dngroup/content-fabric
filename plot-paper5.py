from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from matplotlib import cm


fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

data = pd.DataFrame.from_csv("data.csv")
data = data[data["arrival_time"] == 0.5]

X, Y = np.meshgrid(sorted(list(set(data["te_count"].values))), sorted(list(set(data["client_count"].values))))
Z = np.empty((len(X), len(Y)))


def Zvalue(a, b):
    v=data[(data["te_count"] == a) & (data["client_count"] == b)]["min"].values
    if len(v)==0:
        return np.nan
    else:
        return v[0]

colortuple = ('y', 'b')
colors = np.empty(X.shape, dtype=str)
for y in range(len(Y)):
    for x in range(len(X)):
        colors[x, y] = colortuple[(x + y) % len(colortuple)]


Zvalue = np.vectorize(Zvalue)
Z=Zvalue(X,Y)
surf = ax.plot_wireframe(X, Y, Z,  linewidth=1)
ax.set_xlabel('TE/CP Count')
ax.set_ylabel('client count')
ax.set_zlabel('contract convergeance time')

plt.show()


