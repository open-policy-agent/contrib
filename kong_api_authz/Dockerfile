FROM kong:2.0

LABEL org.wada.name="openpolicyagent/demo-kong-plugin-opa"
LABEL org.wada.description="Kong image with kong-plugin-opa installed"

WORKDIR /usr/kong/opa

# copy the plugin sources
COPY . .

# switch to root to install rocks in /usr/local
USER root

# build and install the plugin
RUN luarocks make

# back to kong user
USER kong