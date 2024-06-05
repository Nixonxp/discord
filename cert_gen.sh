#!/bin/sh

rm ./gateway/cert/*

# 1. Generate CA's private key and self-signed certificate

# -subj:
# /C=RU is for Country
# /ST=Moscow State or province
# /L=Moscow is for Locality name or city
# /O=Balun Cources
# /OU=Education is for Organisation Unit
# /CN=*.balun.courses is for Common Name or domain name
# /emailAddress=leolegrand1014@gmail.com is for email address

openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ./gateway/cert/ca-key.key -out ./gateway/cert/ca-cert.crt -subj "/C=/ST=/L=/O=/OU=/CN=*.ru/emailAddress="

echo "CA's self-signed certificate"
openssl x509 -in ./gateway/cert/ca-cert.crt -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./gateway/cert/server-key.key -out ./gateway/cert/server-req.csr -subj "/C=/ST=/L=/O=/OU=/CN=*.ru/emailAddress="

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
echo "subjectAltName=IP:0.0.0.0,IP:127.0.0.1" > ./gateway/cert/server-ext.cnf
openssl x509 -req -in ./gateway/cert/server-req.csr -days 60 -CA ./gateway/cert/ca-cert.crt -CAkey ./gateway/cert/ca-key.key -CAcreateserial -out ./gateway/cert/server-cert.crt -extfile ./gateway/cert/server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in ./gateway/cert/server-cert.crt -noout -text

# 4. Verify a certificate
openssl verify -CAfile ./gateway/cert/ca-cert.crt ./gateway/cert/server-cert.crt

# 5. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./gateway/cert/client-key.key -out ./gateway/cert/client-req.csr -subj "/C=/ST=/L=/O=/OU=/CN=*.ru/emailAddress="

# 6. Use CA's private key to sign client's CSR and get back the signed certificate
echo "subjectAltName=IP:0.0.0.0,IP:127.0.0.1" > ./gateway/cert/client-ext.cnf
openssl x509 -req -in ./gateway/cert/client-req.csr -days 60 -CA ./gateway/cert/ca-cert.crt -CAkey ./gateway/cert/ca-key.key -CAcreateserial -out ./gateway/cert/client-cert.crt -extfile ./gateway/cert/client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in ./gateway/cert/client-cert.crt -noout -text