local record = function(zone_name, name, content, type='A', proxied=true, ttl=1) {
    zone_name: zone_name,
    name: if std.endsWith(name, zone_name) then name else name + '.' + zone_name,
    content: content,
    type: type,
    proxied: proxied,
    ttl: ttl
};

{
    zone(zone_name):: {
        a(name, ip, proxied=true, ttl=1):: record(zone_name, name, ip, type='A', proxied=proxied, ttl=ttl),
        cname(name, domain, proxied=true, ttl=1):: record(zone_name, name, domain, type='CNAME', proxied=proxied, ttl=ttl),
        txt(name, content):: record(zone_name, name, content, type='TXT', proxied=false, ttl=1),
        mx(name, content):: record(zone_name, name, content, type='MX', proxied=false, ttl=1),
    }
}
