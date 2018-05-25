#! /bin/sh
# Make a self-signed TLS key.
# Thanks to https://github.com/jcbsmpsn/golang-https-example
set -x
openssl req \
	-x509 \
	-nodes \
	-newkey rsa:2048 \
	-keyout .tls.key \
	-out .tls.crt \
	-days 3650 \
	-subj "/C=US/ST=California/L=Willits/O=TPC/OU=TLS/CN=*"
