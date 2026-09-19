# this will install and run knotfree coredns on the dedicated server at vultr

# note that the restarter has moved to AWS. see restarter.py

# todo: - dig for this: dig +short ns1.knotfree.io
export TARGET=$(dig +short ns1.knotfree.io)   # 149.28.250.163

# get token from file. This is a secret file that is not in the repo
export TOKEN=$(cat ~/atw_private/giantToken.txt)
# refill token with func TestMakeGiantTokenToFile() in knotfreeiot

echo $TOKEN

#copy to target, two places.

scp ~/atw_private/giantToken.txt root@$TARGET:/root/giantToken.txt
ssh root@$TARGET 'mkdir atw' 
scp ~/atw_private/giantToken.txt root@$TARGET:/root/atw/giantToken.txt
scp ~/atw/privateKeys4.txt root@$TARGET:/root/atw/privateKeys4.txt


docker build  --platform linux/amd64 -t docker.io/alanwootton2/knotfreecoredns .
docker push docker.io/alanwootton2/knotfreecoredns 

# log in:
# ssh root@$TARGET

ssh root@$TARGET 'ls -lah'

ssh root@$TARGET 'docker pull docker.io/alanwootton2/knotfreecoredns'
 
ssh root@$TARGET 'docker stop $(docker ps -q)' # stop any running containers

# the restarter will do it: ssh root@$TARGET 'docker run -d -e KNOTFREE_TOKEN=$(cat ~/giantToken.txt) -p 53:53/udp -p 53:53/tcp docker.io/alanwootton2/knotfreecoredns  ./coredns'

scp ./restarter.sh root@$TARGET:/root/restarter.sh
ssh root@$TARGET 'ps -ef | grep restarter.sh' # check if it's running. If it is, stop it.
ssh root@$TARGET 'pkill -f restarter.sh' # stop any running restarter
ssh root@$TARGET '/root/restarter.sh  >> /root/restarter.log 2>&1 &'

# then we can check the logs with: ssh root@$TARGET 'tail -f /root/restarter.log'

# ssh root@$TARGET 'docker logs $(docker ps -q) -f'

# todo: start docker on reboot
# todo: re-start docker on fail

