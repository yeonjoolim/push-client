#!/bin/bash

docker login registry.gitlab.example.com
docker push $1

if [ $? -eq 0 ];then
    echo "-----------Success Push Image-----------"
    ./layer-sign $1
    if [ $? -eq 0 ];then
    	echo "-----------Success Sign Image-----------"
    	./tls-client-push $1
    	if [ $? -eq 0 ];then
        	echo "-----------Success Upload Image-----------"
    	else
        	echo "-----------Fail Upload Image-----------"
     	fi
     else
    	echo "-----------Fail Sign Image-----------"
     fi
else                   
    echo "-----------Fail Push Image-----------"
fi
