# cafe
CLI for manage cloudflare records by jsonnet


### How to?

* Install the CLI

```
go install github.com/duythinht/cafe/cmd/cafe
```

* Configuration by those environment
  * CLOUDFLARE_API_TOKEN: token, get from cloudflare
  * ZONES_DIR: folder of zones, contains jsonnet files of records definition, default is `./zones`
  * CAFE_CONFIRM: yes/no, apply to cf (default no, which just dry run)

* RUN, example:

```
♺ duythinht[at]factory ♺ ~/workspace/github.com/duythinht/cafe ♺  ❯
❯❯❯ CAFE_CONFIRM=yes go run cmd/cafe/main.go
Those records will be deleted:
ZONE        TYPE    TTL   NAME                    CONTENT
0x7e6.com   A       1     hello.0x7e6.com         104.21.23.167
0x7e6.com   TXT     1     just-txt.0x7e6.com      hello-world!

Those records will be created:
ZONE        TYPE    TTL   NAME                    CONTENT
0x7e6.com   A       1     hello.0x7e6.com         104.21.13.167
0x7e6.com   TXT     1     just-txt.0x7e6.com      ok-hello-world!
deleting 0x7e6.com   A       1     hello.0x7e6.com         104.21.23.167...
deleting 0x7e6.com   TXT     1     just-txt.0x7e6.com      hello-world!...
creating 0x7e6.com   A       1     hello.0x7e6.com         104.21.13.167... true
creating 0x7e6.com   TXT     1     just-txt.0x7e6.com      ok-hello-world!... true
```
