from PIL import Image
import os
import sys
from PIL import ImageFile
ImageFile.LOAD_TRUNCATED_IMAGES = True
import shutil

def get_size(file):
    # 获取文件大小:KB
    size = os.path.getsize(file)
    return size / 1024

def get_outfile(infile, outfile):
    if outfile:
        return outfile
    dir, suffix = os.path.splitext(infile)
    outfile = '{}-out{}'.format(dir, suffix)
    return outfile

def get_outfile_front(infile, outfile):
    if outfile:
        return outfile
    dir, suffix = os.path.splitext(infile)
    outfile = '{}-front{}'.format(dir, suffix)
    return outfile

def compress_image(infile, outfile='', mb=100, step=10, quality=80):
    """不改变图片尺寸压缩到指定大小
    :param infile: 压缩源文件
    :param outfile: 压缩文件保存地址
    :param mb: 压缩目标，KB
    :param step: 每次调整的压缩比率
    :param quality: 初始压缩比率
    :return: 压缩文件地址，压缩文件大小
    """
    print('Start compress ', infile)
    outfile = get_outfile_front(infile, outfile)
    o_size = get_size(infile)
    if o_size <= mb:
        im = Image.open(infile)
        im.save(outfile, quality=quality) 
        return infile
    while o_size > mb:
        im = Image.open(infile)
        im.save(outfile, quality=quality)
        if quality - step < 0:
            break
        quality -= step
        o_size = get_size(outfile)
    return outfile, get_size(outfile)

#有拉伸效果
def resize_image(infile, outfile='', x_s=480):
    """修改图片尺寸
    :param infile: 图片源文件
    :param outfile: 重设尺寸文件保存地址
    :param x_s: 设置的宽度
    :return:
    """
    im = Image.open(infile)
    x, y = im.size
    #y_s = int(y * x_s / x)
    y_s = 640
    out = im.resize((x_s, y_s), Image.ANTIALIAS)
    outfile = get_outfile(infile, outfile)
    out.save(outfile)


if __name__ == '__main__':
    path = sys.argv[1]
    print(path)
    lsdir = os.listdir(path)
    for i in lsdir:
        compress_image(os.path.join(path, i))
