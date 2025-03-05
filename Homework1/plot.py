import numpy as np
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import matplotlib.colors as mcolors
import matplotlib


X_MIN = 2
X_MAX = 4

def f(x):
    return np.cos(np.exp(x)) / np.sin(np.log(x))

x_values = np.linspace(X_MIN, X_MAX, 400)
y_values = f(x_values)

file_path = r"c:/Users/a.chernyakov/Desktop/BioAlg/Homework1/results.txt"
with open(file_path, "r") as file:
    x_points = list(dict.fromkeys(float(line.strip()) for line in file))

y_points = f(np.array(x_points))

matplotlib.use('Agg')
norm = mcolors.Normalize(vmin=0, vmax=len(x_points) - 1)
cmap = plt.colormaps.get_cmap("RdYlGn")
colors = [cmap(norm(i)) for i in range(len(x_points))]

plt.figure(figsize=(8, 5))
plt.plot(x_values, y_values, label="f(x)", color="blue")

for i, (x, y) in enumerate(zip(x_points, y_points)):
    plt.scatter(x, y, color=colors[i], edgecolor='black', zorder=3)
    #plt.text(x, y, f"{y:.2f}", fontsize=9, ha='left', va='bottom')

plt.xlabel("x")
plt.ylabel("f(x)")
plt.title("График функции с отмеченными точками")
plt.legend()
plt.grid(True)
plt.savefig("Homework1/plot.png", dpi=300)