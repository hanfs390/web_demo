import os
import shutil
import sys

def check(video, picture):
    files = os.listdir(video)
    for i in files:
        if os.path.isdir(os.path.join(video, i)):
            check(os.path.join(video, i), picture)
        else:
            if i.endswith('.mp4'):
                p = i.replace('.mp4', '.jpg')
                if not os.path.exists(os.path.join(picture, p)):
                    print(os.path.join(video, i))
                #else:
                    #shutil.copyfile(os.path.join(picture, p), os.path.join('/home/hfs/pictures', p))
#            else:
#                print(i)

videoDir = sys.argv[1]
pictureDir = sys.argv[2]
print(videoDir, pictureDir)

check(videoDir, pictureDir)
