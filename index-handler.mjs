import dns from 'dns'

console.log('Loading function');

export const handler = async (event, context, callback) => {

    // Node.js program to demonstrate the   
    // dns.resolve() method

    try {
        console.log('Loading function');

        // Set the rrtype for dns.resolve() method
        const rrtype = "A";

        // Calling dns.resolve() method for hostname
        // geeksforgeeks.org and print them in
        // console as a callback
        // let resolver = new dns.Resolver();
        // resolver.setServers(['149.28.250.163']);
        console.log('calling setServers');

        dns.setServers(['149.28.250.163']);
        let site = 'a-person-channel.vr';
        // site = 'google.com';

        console.log('calling resolve');

        dns.resolve(site, rrtype, (err, records) => {
            console.log('IN resolve');
            if (err) {
                console.log('dns.resolve error: %j', err);
                callback(err, '');
                // throw err
            }
            console.log('records records records records: %j', records);
            callback(null, records[0]);
        });
        //console.log('Received event:', JSON.stringify(event, null, 2));
        console.log('event =', event);
        // console.log('value2 =', event.key2);
        // console.log('value3 =', event.key3);

        // return event.key1;  // Echo back the first key value
        // throw new Error('Something went wrong');
    } catch (e) {
        console.log('error: %j', e);
        callback(e, '');
        throw e
    }
};
