#bash -ef

# this is because the docker logs are filling up the disk and causing the server to crash. 
# This restarter will stop the docker process, clean the logs, and start it again. 
# It will do this every day. 
# This is a temporary solution until I can figure out a better way to manage the logs.

# for always: 
# sleep
# stop the docker process
# clean the logs
# start again

echo restarter starting loop at $(date)

docker stop $(docker ps -q)

docker run -d -e KNOTFREE_TOKEN=$(cat ~/giantToken.txt) -p 53:53/udp -p 53:53/tcp docker.io/alanwootton2/knotfreecoredns  ./coredns

echo restarter running now $(docker ps -q) at $(date)

while true; do
    echo sleeping 12 hours $(date) # 60 # 86400 is one day  12h is 43200
    # sleep 12h - check docker.ps every 10 seconds if something's running for 12 hours
    elapsed=0
    max_time=43200  # 12 hours in seconds
    
    while [ $elapsed -lt $max_time ]; do
        if docker ps -q | grep -q .; then
            # too boring: echo docker running: $(docker ps -q) at $(date)
            sleep 1
        else
            echo docker not running at $(date)
            # its dead Jim. Let's break out of the sleep loop and restart it.
            break
        fi
        sleep 10
        elapsed=$((elapsed + 10))
    done

    echo stopping $(date)
    docker stop $(docker ps -q)
    echo remove old containers $(date)
    # docker rm -v -f $(docker ps -qa) ???? 
    docker system prune --force
    echo start coredns docker again $(date)
    docker run -d -e KNOTFREE_TOKEN=$(cat ~/giantToken.txt) -p 53:53/udp -p 53:53/tcp docker.io/alanwootton2/knotfreecoredns  ./coredns
    echo restarted docker: $(docker ps -q) at $(date)

done