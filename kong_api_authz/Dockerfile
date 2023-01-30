FROM kong:2.8.3

LABEL org.wada.name="openpolicyagent/demo-kong-plugin-opa"
LABEL org.wada.description="Kong image with kong-plugin-opa installed"

WORKDIR /usr/kong/opa

# copy the plugin sources
COPY . .

# switch to root to install rocks in /usr/local
USER root

# install dependencies
RUN apk update && \
  apk add gcc libc-dev && \
  luarocks install lua-cjson && \
  luarocks install lua-resty-http && \
  luarocks install lua-resty-jwt

# build and install the plugin
RUN luarocks make

# back to kong user
USER kong
