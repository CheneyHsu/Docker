#!/bin/bash

cd ~; mkdir aufs2;cd aufs2
mkdir container image-layer{1..4} mnt
for i in container image-layer{1..4} ; do echo "i am $i" > $i/$i.txt;done
sudo mount -v -t aufs -o dirs=./container:./image-layer4:./image-layer3:./image-layer2:./image-layer1 none ./mnt
sudo cat /sys/fs/aufs/si*/*
echo -e "\n write to mnt images-layer1.txt" >>mnt/image-layer4.txt
cat mnt/image-layer4.txt
cat image-layer4/image-layer4.txt
ls container/
~