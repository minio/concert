# concert
<blockquote>
To use this tool you need to have a valid domain name. This tool doesn't support acme DNS challenges.
</blockquote>

Concert is a console based certificate generating tool using letsencrypt, built using really simple [ACME library](https://github.com/xenolf/lego).

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

Generate certificates in custom directory.
```
$ sudo concert gen --dir my-certs-dir <EMAIL> <DOMAIN>
```

Renew certificates in `certs` directory by default.
```
$ sudo concert renew <EMAIL>
```

Generate certificates in custom directory.
```
$ sudo concert renew --dir my-certs-dir <EMAIL>
```

<blockquote>
Concert requires root on all platforms other than Linux. On linux you can use
setcap permissions to allow concert to have permissions to bind ports.
</blockquote>

```
setcap cap_net_bind_service=+ep `which concert`
```
