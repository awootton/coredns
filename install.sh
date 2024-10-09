# this will install and run knotfree coredns on the dedicated server at vultr

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


docker build -t gcr.io/fair-theater-238820/knotfreecoredns .
docker push gcr.io/fair-theater-238820/knotfreecoredns 

# log in:
# ssh root@$TARGET

ssh root@$TARGET 'ls -lah'

ssh root@$TARGET 'docker pull gcr.io/fair-theater-238820/knotfreecoredns'
 
ssh root@$TARGET 'docker stop $(docker ps -q)' # stop any running containers

ssh root@$TARGET 'docker run -d -e KNOTFREE_TOKEN=$(cat ~/giantToken.txt) -p 53:53/udp -p 53:53/tcp gcr.io/fair-theater-238820/knotfreecoredns  ./coredns'

# ssh root@$TARGET 'docker logs $(docker ps -q) -f'

# todo: start docker on reboot
# todo: re-start docker on fail

