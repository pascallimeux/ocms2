#!/bin/bash
echo "make a tar.gz for ocms2 and send it to the integration server"
sh env.sh
tar cvzf ocms2.tar.gz /data/ocms2/dist
scp ocms2.tar.gz orange@10.194.18.46:/tmp
rm ocms2.tar.gz