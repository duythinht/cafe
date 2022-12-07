# file 0x7e6.com.jsonnet

local cafe = import 'cafe.libsonnet';

local records = cafe.zone('0x7e6.com');

[
    records.a('hello.0x7e6.com', '104.20.13.167', ttl=1),
    records.a('hello.0x7e6.com', '104.20.13.137', ttl=1, proxied=false),
    records.a('test', '104.21.63.157'),
    records.cname('test-cname', 'test.example.com'),
    records.cname('_A4B9EA64EED76C439C5B8C4F21A0EDFD.app.0x7e6.com', '8CD85996A4306A0328226737856BC22B.46A116D49DA29D37FAF3505F74979C0F.f3aa087e591c0f5.comodoca.com', proxied=false, ttl=3600),
    records.txt('just-txt', 'sad')
]
