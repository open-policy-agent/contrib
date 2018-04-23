FROM confluentinc/cp-kafka:4.0.0-3

COPY ./target/kafka-authorizer-opa-VERSION-package/share/java/kafka-authorizer-opa/kafka-authorizer-opa-VERSION.jar /usr/share/java/kafka/
COPY ./target/kafka-authorizer-opa-VERSION-package/share/java/kafka-authorizer-opa/gson-2.8.2.jar /usr/share/java/kafka/
COPY ./secrets /etc/kafka/secrets