FROM ubuntu
RUN mkdir opa
WORKDIR opa
COPY bin/opa-mongo .
RUN chmod +x opa-mongo
CMD ["./opa-mongo", "run", "-c", "/opa/config.yaml", "-p", "/opa/example.rego"]