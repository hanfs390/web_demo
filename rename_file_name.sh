#!/bin/sh
#change file

rename_all(){
	echo "第一个参数为 $1 !"
	rename "s/ /_/g" $1/*
	rename "s/yishesp.com/_/g" $1/*
	rename "s/_//g" $1/*
	for file in $1/*; do
    		if [ -d $file ];then
		echo $file
		rename_all $file
    		fi
	done
}

rename_all "/home/hfs/Desktop/ke/txt"
