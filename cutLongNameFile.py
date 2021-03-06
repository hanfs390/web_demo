import os
def cutName(path):
    files = os.listdir(path)
    for i in files:
        if os.path.isdir(os.path.join(path, i)):
            cutName(os.path.join(path, i))
        else:
            if len(i) > 40:
                new = i[20:]
                os.rename(os.path.join(path, i), os.path.join(path, new)) 
                print(new)
            else:
                print(i)


path=input('请输入文件路径：')       
cutName(path)
