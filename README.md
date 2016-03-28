# Concert
Concert is a console based certificate generation tool for [letsencrypt.org](https://letsencrypt.org/). `Let’s Encrypt` is a free (as in free beer), automated, and open certificate authority.

## Prerequisite
* A valid domain name purchased from any domain registrar.
* `root` access to the server pointed by the domain name.
* Working email address for the domain.

### Download
We **STRONGLY RECOMMEND** installing `concert` from source, because it requires root access.

### Compile from Source (RECOMMENDED)
Assuming that you golang installed already, run the following command to download and install `concert` from source.
```
$ go get -u github.com/minio/concert
```
## How to generate a certificate?
To generate a certificate and key for `example.com` under `my-certs` directory.
```bash
sudo concert gen --dir my-certs admin@example.com example.com
sudo ls my-certs
certs.json public.crt private.key
```
NOTE: Generated certificates are valid only for a maximum of 90 days. Please visit the following link for more details - [https://letsencrypt.org/2015/11/09/why-90-days.html](https://letsencrypt.org/2015/11/09/why-90-days.html)

## How to renew a certificate?
To renew a certificate for example.com under ‘certs’ directory. New certs are generated and saved in the same directory as before.
```bash
sudo concert renew admin@example.com
```

## How to automatically renew certificates?
You can run `concert` in server mode to automatically renew certificates, once in every 45 days.
```bash
sudo concert server --dir my-certs admin@example.com example.com
```

### FAQ
* Why `concert` requires root access?
ACME protocol requires root access to verify authenticity of the domain ownership. During the certification generation phase, `concert` temporarily listens on port `80` or `443` to allow letsencrypt.org service connect and verify the ownership. Only `root` is allowed to bind to any port below `1024`.

* Can I run `concert` as non-root?
On GNU/Linux, it is possible to run as non-root by granting bind only access to  `concert`.
```
sudo setcap cap_net_bind_service=+ep `which concert`
```
