local cafe = import 'cafe.libsonnet';
local records = cafe.zone('0x97a.com');

[
    records.a('radio.0x97a.com', '66.241.124.85', ttl=1, proxied=false),
]
