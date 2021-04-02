#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

export LD_LIBRARY_PATH=$(dirname $PWD)/:$LD_LIBRARY_PATH
export PATH=$(dirname $PWD)/prebuilt/linux:$(dirname $PWD)/prebuilt/win64:$PATH
export WASMER_BACKTRACE=1
cd ..

pid=`ps -ef | grep chainmaker | grep "\-c ../config/wx-org1/chainmaker.yml local-tbft" | grep -v grep |  awk  '{print $2}'`
if [ -z ${pid} ];then
    nohup ./chainmaker start -c ../config/wx-org1/chainmaker.yml local-tbft > panic.log &
    echo "wx-org1 chainmaker is startting, pls check log..."
else
    echo "wx-org1 chainmaker is already started"
fi

pid2=`ps -ef | grep chainmaker | grep "\-c ../config/wx-org2/chainmaker.yml local-tbft" | grep -v grep |  awk  '{print $2}'`
if [ -z ${pid2} ];then
    nohup ./chainmaker start -c ../config/wx-org2/chainmaker.yml local-tbft > panic.log &
    echo "wx-org2 chainmaker is startting, pls check log..."
else
    echo "wx-org2 chainmaker is already started"
fi



pid3=`ps -ef | grep chainmaker | grep "\-c ../config/wx-org3/chainmaker.yml local-tbft" | grep -v grep |  awk  '{print $2}'`
if [ -z ${pid3} ];then
    nohup ./chainmaker start -c ../config/wx-org3/chainmaker.yml local-tbft > panic.log &
    echo "wx-org3 chainmaker is startting, pls check log..."
else
    echo "wx-org3 chainmaker is already started"
fi


pid4=`ps -ef | grep chainmaker | grep "\-c ../config/wx-org4/chainmaker.yml local-tbft" | grep -v grep |  awk  '{print $2}'`
if [ -z ${pid4} ];then
    nohup ./chainmaker start -c ../config/wx-org4/chainmaker.yml local-tbft > panic.log &
    echo "wx-org4 chainmaker is startting, pls check log..."
else
    echo "wx-org4 chainmaker is already started"
fi

# nohup ./chainmaker start -c ../config/wx-org5/chainmaker.yml local-tbft > panic.log &