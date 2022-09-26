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
