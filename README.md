# RinGo

## Abstract
With a fleet of 100's of cams streaming data, it would be less than ideal for each cam to have its own DB connection to persist the data. This is where RinGo comes into play. Its a light weight middle man between a fleet of IOT devices producing data, a ML inference endpoint, and a persistance layer. RinGo will accept data from IOT devices via http, hold in a ring buffer (hence the ring in RingGo) and periodically send this data to inference endpoint. Depending on the response, data will be persisted. This keeps a nice separation between IOT device, ML layer, and persistance layer. 
