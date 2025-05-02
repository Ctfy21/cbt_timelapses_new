#!bin/bash

export LANG=""; 

SAMBA_IP="//192.168.1.1/snapshots"
HOST_IP=$(hostname -I)

watch -n 86400 mount.cifs $SAMBA_IP /mnt/seefetch -o user=sambauser,pass=111111111

touch /home/cbt_timelapses_new/cbt_timelapses_frontend/.env 
echo "VUE_APP_PATH_START=${HOST_IP}" >> /home/cbt_timelapses_new/cbt_timelapses_frontend/.env

redis-server &&
serve -p 8080 -s /home/cbt_timelapses_new/cbt_timelapses_frontend/dist/ &&
/home/cbt_timelapses_new/cbt_timelapses_backend/main
