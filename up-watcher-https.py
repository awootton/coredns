import subprocess
import time
import json

try:
    import requests
except ImportError:
    print("Error: requests library not installed. Install it with: pip install requests")
    exit(1)

# I run this when I doubt the reliability of the system. 
# sometimes shit breaks. 
# note that this is not testing coredns.

# here's the schema
# //   "Status": 0,
# //   "TC": false,
# //   "RD": true,
# //   "RA": true,
# //   "AD": true,
# //   "CD": false,
# //   "Question": [
# //     {
# //       "name": "example.com",
# //       "type": 16
# //     }
# //   ],
# //   "Answer": [
# //     {
# //       "name": "example.com",
# //       "type": 16,
# //       "TTL": 300,
# //       "data": "\"v=spf1 -all\""
# //     },
# //     {
# //       "name": "example.com",
# //       "type": 16,
# //       "TTL": 300,
# //       "data": "\"_k2n1y4vw3qtb4skdx9e7dxt97qrmmq9\""
# //     }
# //   ]
# // }


domain = "meta_group_id.testmain-0n0u0e16p-0.vr"  # Replace with your target domain

count = 0
goodcount = 0
failcount = 0
while True:
    # also http to knotfree.com:8085 for local test.
    url = "https://knotfree.net/api1/dns-query?name=meta_group_id.testmain-0n0u0e16p-0.vr&type=TXT&knotfree=1"
    
    # Get the URL
    
    expectedText =  "meta_group_id-no-leading-underscore"
    receivedText = ""
    for i in range(10):
        try:
            val = requests.get(url)
            if val.status_code != 200:
                print(f"Error getting URL: {val.status_code}")
                continue
            json_string = val.text
            data = json.loads(json_string)
            if data["Status"] != 0:
                print(f"Error in DNS response: Status {data['Status']}")
                # try again
                continue
            receivedText = data["Answer"][0]["data"] if data.get("Answer") else ""
            # marshal receivedText into JSON 
            # if receivedJson.Status 
            print(f"Received text: {receivedText} at {time.strftime('%X')}")
            break
        except Exception as e:
            print(f"Error getting URL: {e}")
            
    if receivedText != expectedText:
        print(f"FAIL FAIL FAIL Expected '{expectedText}' but got '{receivedText}' at {time.strftime('%X')}")
        failcount +=1
    else:
        goodcount +=1
    
            
    # Wait 10 seconds before executing the next dig
    time.sleep(10)
    count += 1
    if count % 10 == 0:
        percent = (goodcount / count) * 100 if count > 0 else 0
        print(f"--- {count} iterations, {goodcount} good, {failcount} fails ({percent:.2f}%) at {time.strftime('%X')} ---")

