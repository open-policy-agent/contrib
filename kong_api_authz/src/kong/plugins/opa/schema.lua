local typedefs = require "kong.db.schema.typedefs"

return {
  name = "opa",
  fields = {
    {
      -- this plugin will only be applied to Services or Routes
      consumer = typedefs.no_consumer
    },
    {
      -- this plugin will only run within Nginx HTTP module
      protocols = typedefs.protocols_http
    },
    {
      config = {
        type = "record",
        fields = {
          -- Plugin's configuration's schema
          {
            server = {
              type = "record",
              fields = {
                {
                  protocol = typedefs.protocol {
                    default = "http"
                  },
                },
                {
                  host = typedefs.host {
                    default = "localhost"
                  },
                },
                {
                  port = {
                    type = "number",
                    default = 8181,
                    between = { 0, 65534 },
                  },
                },
                {
                  connection = {
                    type = "record",
                    fields = {
                      {
                        timeout = {
                          type = "number",
                          default = 60,
                        },
                      },
                      {
                        pool = {
                          type = "number",
                          default = 10,
                        },
                      },
                      {
                        read_timeout = {
                          type = "number",
                          default = 1000,
                        },
                      },
                      {
                        send_timeout = {
                          type = "number",
                          default = 1000,
                        },
                      },
                      {
                        connect_timeout = {
                          type = "number",
                          default = 1000,
                        },
                      }
                    },
                  },
                },
              },
            },
          },
          {
            document = {
              type = "record",
              fields = {
                {
                  include_headers = {
                    type = "array",
                    elements = {
                      type = "string"
                    },
                  },
                },
              },
            },
          },
          {
            policy = {
              type = "record",
              fields = {
                {
                  base_path = {
                    type = "string",
                    default = "v1/data"
                  },
                },
                {
                  decision = {
                    type = "string",
                    required = true
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}
