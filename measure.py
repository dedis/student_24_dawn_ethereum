import numpy as np

def main():
    arr10 = np.array([9.0728, 13.0841, 9.4433, 7.7694, 6.913, 10.0022, 9.3505, 9.5906, 6.6857, 8.7748])
    print("batch 10 mean: {}, std: {} in ms".format(arr10.mean(), arr10.std()))
    
    arr50 = np.array([53.6601, 57.2375, 56.673, 51.3562, 53.0928, 38.9681, 37.7452, 35.4205, 37.6897, 38.2749])
    print("batch 50 mean: {}, std: {} in ms".format(arr50.mean(), arr50.std()))
    
    arr100 = np.array([115.1835, 58.3309, 74.9841, 122.2076, 71.8468, 71.5715, 125.4863, 74.1998, 106.3774, 115.4975])
    print("batch 100 mean: {}, std: {} in ms".format(arr100.mean(), arr100.std()))
    
def measure():
    file_name = "exelogfile"
    with open(file_name, "r") as f:
        lines = f.readlines()
        prefix = "[ENC][EXE][Elapse]: "
        print(lines[0][len(prefix):-3])
        modified = list(map(lambda x: float(x[len(prefix):-3]), lines))
        arr = np.array(modified)
        print("exe: mean: {}, std: {}".format(np.mean(arr), np.std(arr)))
    
if __name__ == "__main__":
    measure()

# fig one: batch transaction execution time, with 3 f3b nodes, 10 trails
# 1: exe: mean: 119.02683000000002, std: 9.01195281590511
# 5: exe: mean: 538.6084222222222, std: 40.35938004079453
# 10: exe: mean: 1147.3095, std: 28.542923645800515
# 20: exe: mean: 2217.4526, std: 18.558086039847968
# 50: exe: mean: 5877.7835000000005, std: 44.25313713543366


# fig two: batch transactions=10, with 3, 5, 10, 20 f3b nodes, 5 trails, share size
# 3: exe: mean: 1147.3095, std: 28.542923645800515    886
# 5: exe: mean: 1381.59506, std: 37.69252243984744    1460
# 10: exe: mean: 1750.23148, std: 74.13925567566488   2895
# 20: exe: mean: 2834.56372, std: 318.7618507839946   5775
# 40: exe: mean: 6282.05756, std: 909.1281258529376   11535 

# share&proof  size bytes: 287n+25
# 35KB for n=128


