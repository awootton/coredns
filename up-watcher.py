import subprocess
import time

# I run this when I doubt the reliability of the system. 
# sometimes shit breaks. 

domain = "meta_group_id.testmain-0n0u0e16p-0.vr"  # Replace with your target domain

count = 0
goodcount = 0
failcount = 0
while True:
    # Run the dig command (capturing output)
    result = subprocess.run(['dig', "@149.28.250.163" ,domain, "TXT"], capture_output=True, text=True)
    
    # Print the timestamp and output to your console
    got = result.stdout
    if not "status: NOERROR" in got:
        print(f"--- DIG fail at {time.strftime('%X')} ---")
        print(result.stdout)
        failcount +=1
    else: # don't print the good ones
        goodcount +=1
    
    # check for error 
    if result.returncode != 0:
        print(f"Error running dig: {result.stderr}")       

    # Extract and print the TXT record value
    found = ""
    for line in result.stdout.split('\n'):
        # print(f"TXT Record: {line}")
        if 'IN TXT' in line:
            # print(f"TXT Record: {line}")
            txt_value = line[line.find('IN TXT'):]
            if len(txt_value) > 0:
                found = txt_value
            
    if found:
        # txt_value = found[found.find('IN TXT'):]
        print(f"Found TXT Record: {found} at {count} when {time.strftime('%X')}")
        if not "meta_group_id-no-leading-underscore" in found:
            print(f"FAIL FAIL FAIL Found fail in TXT record! {time.strftime('%X')}")
    else:
        print(f"FAIL FAIL FAIL No TXT record found. {time.strftime('%X')}")
            
    # Wait 30 seconds before executing the next dig
    time.sleep(30)
    count += 1
    if count % 10 == 0:
        percent = (goodcount / count) * 100 if count > 0 else 0
        print(f"--- {count} iterations, {goodcount} good, {failcount} fails ({percent:.2f}%) at {time.strftime('%X')} ---")

# note the status. Not the same as not found. 
# (1 server found)
# ;; global options: +cmd
# ;; Got answer:
# ;; ->>HEADER<<- opcode: QUERY, status: SERVFAIL, id: 2716
# ;; flags: qr rd; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
# ;; WARNING: recursion requested but not available

# ;; OPT PSEUDOSECTION:
# ; EDNS: version: 0, flags:; udp: 4096
# ;; QUESTION SECTION:
# ;meta_group_id.testmain-0n0u0e16p-0.vr. IN TXT

# ;; Query time: 2068 msec
# ;; SERVER: 149.28.250.163#53(149.28.250.163)
# ;; WHEN: Thu Jun 04 18:09:50 PDT 2026
# ;; MSG SIZE  rcvd: 66


# THIS is a not found:

# ; <<>> DiG 9.10.6 <<>> @149.28.250.163 meta_group_id.Xtestmain-0n0u0e16p-0.vr TXT
# ; (1 server found)
# ;; global options: +cmd
# ;; Got answer:
# ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 11609
# ;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
# ;; WARNING: recursion requested but not available

# ;; OPT PSEUDOSECTION:
# ; EDNS: version: 0, flags:; udp: 4096
# ;; QUESTION SECTION:
# ;meta_group_id.Xtestmain-0n0u0e16p-0.vr.        IN TXT

# ;; ANSWER SECTION:
# meta_group_id.xtestmain-0n0u0e16p-0.vr. 0 IN TXT "error: topic not found"

# ;; Query time: 127 msec
# ;; SERVER: 149.28.250.163#53(149.28.250.163)
# ;; WHEN: Thu Jun 04 18:11:41 PDT 2026
# ;; MSG SIZE  rcvd: 140
