# file 0x7e6.com.jsonnet

local cafe = import 'cafe.libsonnet';

local records = cafe.zone('0x97a.com');

[
    records.a('0x97a.com', '66.241.125.167', ttl=1, proxied=false),
    records.a('testing.0x97a.com', '34.160.220.95', ttl=1, proxied=false),
    records.deprecated('testing.0x97a.com'),
]
