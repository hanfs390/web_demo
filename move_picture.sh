#!/bin/sh
#change file

rename_all(){
	for file in $1/*; do
		mv $file/* /home/hfs/others/picture/欧美色图/综合	
	done
}

rename_all "/home/hfs/go/src/go_spider/images"
