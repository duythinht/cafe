# file 0x7e6.com.jsonnet

local cafe = import 'cafe.libsonnet';

local records = cafe.zone('0x7e6.com');

[
    records.mx('0x7e6.com', 'alt4.gmr-smtp-in.l.google.com', priority=40),
    records.mx('0x7e6.com', 'alt3.gmr-smtp-in.l.google.com', priority=30),
    records.mx('0x7e6.com', 'alt2.gmr-smtp-in.l.google.com', priority=20),
    records.mx('0x7e6.com', 'alt1.gmr-smtp-in.l.google.com', priority=10),
    records.mx('0x7e6.com', 'gmr-smtp-in.l.google.com', priority=5)
]