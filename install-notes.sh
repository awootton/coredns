
# this will install and run knotfree coredns on the dedicated server at vultr

export TARGET=149.28.250.163

# get token from file. This is a secret file that is not in the repo
export TOKEN=$(cat ~/atw_private/giantToken.txt)
# refil token with func TestMakeGiantTokenToFile() in knotfreeiot

echo $TOKEN

#copy to target

scp ~/atw_private/giantToken.txt root@$TARGET:/root/giantToken.txt
ssh root@$TARGET 'mkdir atw' 
scp ~/atw/giantToken.txt root@$TARGET:/root/atw/giantToken.txt
scp ~/atw/privateKeys4.txt root@$TARGET:/root/atw/privateKeys4.txt


# log in:
ssh root@$TARGET

ssh root@$TARGET 'ls -lah'


# did these by hand:

sudo apt update

apt-get install docker.io

docker --version

docker pull gcr.io/fair-theater-238820/knotfreecoredns

export KNOTFREE_TOKEN=$(cat ~/giantToken.txt)
echo $KNOTFREE_TOKEN

# remove all containers
# how else to get rid of all the old logs? 
docker stop $(docker ps -q)
docker rm -v -f $(docker ps -qa)

# todo: set this up to start on boot
# docker run  -e KNOTFREE_TOKEN=$KNOTFREE_TOKEN -p 53:53/udp -p 53:53/tcp gcr.io/fair-theater-238820/knotfreecoredns  ./coredns 
 
docker run -d -e KNOTFREE_TOKEN=$(cat ~/giantToken.txt) -p 53:53/udp -p 53:53/tcp gcr.io/fair-theater-238820/knotfreecoredns  ./coredns 

docker logs $(docker ps -q) -f

docker stop $(docker ps -q)



