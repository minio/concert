# concert

Concert is a console based certificate generating tool using letsencrypt.

### Install

You need to have golang installed to compile `concert`.
```
$ go get -u github.com/minio/concert
```

### How to run?

Generates certs in `certs` folder by default.
```
$ concert gen <EMAIL> <DOMAIN>
```


Generate certificate other than the default folder.

```
$ concert gen --folder my-folder <EMAIL> <DOMAIN>
```
