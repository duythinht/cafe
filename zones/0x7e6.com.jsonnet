# file 0x7e6.com.jsonnet

local cafe = import 'cafe.libsonnet';

local records = cafe.zone('0x7e6.com');

[
    records.a('hello.0x7e6.com', '104.20.13.167', ttl=1),
    records.a('hello.0x7e6.com', '104.20.13.137', ttl=1, proxied=false),
    records.a('test', '104.21.63.157'),
    records.cname('test-cname', 'test.example.com'),
    records.txt('just-txt', 'sad')
]
