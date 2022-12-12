from multiprocessing import Process
import os
import glob
import subprocess
import time

def initialize_datadir(num: int):
    for i in range(1, 3):
        cmd = "geth --datadir .ethereum{}/ init ../multiclique.json".format(i)
        print(subprocess.check_output(cmd.split(" ")))
    
    
def delete(id: int):
    if os.path.exists(".\\ethereum{}\\geth".format(id)):
        os.system("rd /s /q .\\.ethereum{}\\geth".format(id))
    

def main(num: int):
    # remove everything + initialize node
    for i in range(1, num+1):
        delete(i)
        # pass
        
    initialize_datadir(num)
    

if __name__ == '__main__':
    main(2)