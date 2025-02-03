# file 0x7e6.com.jsonnet

local cafe = import 'cafe.libsonnet';

local records = cafe.zone('0x97a.com');

[
    records.a('*', '66.241.125.167', ttl=1, proxied=false)
]
