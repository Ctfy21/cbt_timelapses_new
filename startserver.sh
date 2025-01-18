#!bin/bash

export LANG=""; 

IP=$(hostname -I)

touch .env 
echo "LOCAL_IP=${IP}" >> .env
redis-server &&
serve -p 8080 -s /home/blunder/bin/cbt_timelapses_new/cbt_timelapses_frontend/dist/ &&
/home/blunder/bin/cbt_timelapses_new/cbt_timelapses_backend/main
