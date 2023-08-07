#%%
import numpy as np
from raster import saveToCOG, makeTransform


# %%

bbox = {"latMin": 0, "lonMin": 0, "latMax": 10, "lonMax": 10}

data = np.zeros((512, 512))
data[:256, :256] = 1
data[256:, :256] = 2
data[256:, 256:] = 3
data[:256, 256:] = 4


transform = makeTransform(512, 512, bbox)

cogFile = saveToCOG("testfile.tiff", data, "EPSG:4326", transform, noDataVal=-9999, mode="direct")
