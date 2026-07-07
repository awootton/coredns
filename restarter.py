import subprocess
import sys
import json


### is tjis? deleteme 

def lambda_handler(event, context):
    
    domain = "a-person-channel.vr"   

    try:
        result = subprocess.run(
            ["dig", "@149.28.250.163",  domain],
            capture_output=True, text=True, check=True
        )

        result = subprocess.run(
            ["dig", "@149.28.250.163",  domain],
            capture_output=True, text=True, check=True
        )
        # ips = [line for line in result.stdout.strip().split('\n') if line and not line.startswith(';')]
        ips = result.stdout # [result.stdout.strip().split('\n')]
        # print(f"IP addresses for {domain}: {ips}")
        # print(f"IP addresses for {domain}: {ips}")

        index = ips.index("216.128.128.195")
        if index > 0:
            print("good exit")
            return {
            'statusCode': 200,
            'body': json.dumps('ok')
        } 
        return {
            'statusCode': 500,
            'body': json.dumps('no stdout from dig')
        }
        
    except subprocess.CalledProcessError as e:
        print(f"Error running dig: {e}", file=sys.stderr)
        return {
            'statusCode': 500,
            'body': json.dumps('exception from dig')
        }

lambda_handler(None, None)

 
import socket

# make an envrroent
# pip3 install dnspython
import dns.resolver

resolver = dns.resolver.Resolver()
resolver.nameservers = ['149.28.250.163']
dns.resolver.override_system_resolver(resolver)

data = "" # socket.gethostbyname("a-person-channel.vr")

try:
    data = socket.gethostbyname("a-person-channel.vr")
except socket.error:
    print("Error: Unable to resolve hostname",socket.error)
    data = "Error: Unable to resolve hostname"
# ip = repr(data)

print(data)