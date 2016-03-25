# concert
<blockquote>
To use this tool you need to have a valid domain name. This tool doesn't support the acme DNS challenges yet.
</blockquote>

Concert is a console based certificate generating tool using letsencrypt.

### Install

You need to have golang installed to compile `concert`.
```
$ go get -u github.com/minio/concert
```

### How to run?

Generates certs in `certs` directory by default.
```
$ sudo concert gen <EMAIL> <DOMAIN>
```

Generate certificate other than the default directory.

```
$ sudo concert gen --dir certs-dir <EMAIL> <DOMAIN>
```

## On linux

Running as root might not be advisable some times. Use setcap instead on linux:
```
setcap cap_net_bind_service=+ep `which concert`
```
