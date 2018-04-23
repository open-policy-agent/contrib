#!/bin/bash

# This script is based on https://github.com/confluentinc/cp-docker-images/blob/master/examples/kafka-cluster-ssl/secrets/create-certs.sh.

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace

PASSWORD="tutorial"

# Generate CA key
openssl req -new -x509 -keyout snakeoil-ca-1.key -out snakeoil-ca-1.crt -days 365 -subj '/CN=ca1.tutorial.openpolicyagent.org/OU=TUTORIAL/O=OPA/L=SF/S=CA/C=US' -passin pass:$PASSWORD -passout pass:$PASSWORD

for i in broker pii_consumer anon_consumer fanout_producer anon_producer
do
    echo $i

    # Create keystores
    keytool -genkey -noprompt \
             -alias $i \
             -dname "CN=$i.tutorial.openpolicyagent.org, OU=TUTORIAL, O=OPA, L=SF, S=CA, C=US" \
             -keystore kafka.$i.keystore.jks \
             -keyalg RSA \
             -storepass $PASSWORD \
             -keypass $PASSWORD

    # Create CSR, sign the key and import back into keystore
    keytool -keystore kafka.$i.keystore.jks -alias $i -certreq -file $i.csr -storepass $PASSWORD -keypass $PASSWORD
    openssl x509 -req -CA snakeoil-ca-1.crt -CAkey snakeoil-ca-1.key -in $i.csr -out $i-ca1-signed.crt -days 9999 -CAcreateserial -passin pass:$PASSWORD
    echo "yes" | keytool -keystore kafka.$i.keystore.jks -alias CARoot -import -file snakeoil-ca-1.crt -storepass $PASSWORD -keypass $PASSWORD
    echo "yes" | keytool -keystore kafka.$i.keystore.jks -alias $i -import -file $i-ca1-signed.crt -storepass $PASSWORD -keypass $PASSWORD

    # Create truststore and import the CA cert.
    echo "yes" | keytool -keystore kafka.$i.truststore.jks -alias CARoot -import -file snakeoil-ca-1.crt -storepass $PASSWORD -keypass $PASSWORD

    echo "$PASSWORD" > ${i}_sslkey_creds
    echo "$PASSWORD" > ${i}_keystore_creds
    echo "$PASSWORD" > ${i}_truststore_creds


    cat > ${i}.ssl.config <<EOF
ssl.truststore.location=/etc/kafka/secrets/kafka.${i}.truststore.jks
ssl.truststore.password=tutorial
ssl.keystore.location=/etc/kafka/secrets/kafka.${i}.keystore.jks
ssl.keystore.password=tutorial
ssl.key.password=tutorial
security.protocol=SSL
EOF
done
