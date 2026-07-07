

import {handler} from './index-handler.mjs';

handler({key1: 'value1', key2: 'value2', key3: 'value3'}, {}, (err, result) => {
    if (err) {
        console.error('Error:', err);
    } else {
        console.log('Result:', result);
    }
}
);
