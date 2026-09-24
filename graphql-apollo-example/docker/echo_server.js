const { ApolloServer, gql, AuthenticationError } = require('apollo-server');
// Need this function to be able to hand a proper GraphQL string to OPA,
// since the `gql` function creates an Object type.
// Ref: https://github.com/apollographql/graphql-tag/issues/144
const { print } = require('graphql');
const axios = require('axios');

// Get environment variable values at startup.
const host = process.env.HOST || "0.0.0.0";
const port = process.env.PORT || "5000";
const opa_url = process.env.OPA_ADDR || "http://localhost:8181";
const policy_path = process.env.POLICY_PATH || "/v1/data/httpapi/authz";


// Resolvers define the technique for fetching the types defined in the
// schema. This resolver retrieves books from the "books" array above.
const resolvers = {
  Query: {
    employeeByID(parent, args, context, info) {
      return {
        id: context.user,
        salary: Math.floor(Math.random() * 10000),
      };
    },
  }
};


// Query OPA for an authz decision, using an input JSON object.
async function opa_authz(target_url, query, schema, user, variables, token) {
  var out = false;
  var input = {
    input: {
      schema: schema,
      query: query,
      user: user,
      variables: variables,
    }
  };
  if (token !== undefined) {
    input.input.token = token; // Add JWT token if available.
  }
  console.log("Checking input...");
  console.log(input);
  await axios
    .post(target_url, input)
    .then(res => {
      out = res.data.result.allow;
      console.log("Auth response:");
      console.log(res.data.result);
    })
    .catch(error => {
      console.error(error);
      out = false;
    });
  return out;
}


// A schema is a collection of type definitions (hence "typeDefs")
// that together define the "shape" of queries that are executed against
// your data.
const typeDefs = gql`
type Employee {
  id: String!
  salary: Int!
}

schema {
  query: Query
}

type Query {
  employeeByID(id: String!): Employee
}
`;


// The ApolloServer constructor requires two parameters: your schema
// definition and your set of resolvers. Other parameters are optional.
const server = new ApolloServer({
  debug: false, // Set to true, or remove in dev.
  typeDefs,
  resolvers,
  csrfPrevention: true,
  context: async ({ req }) => {
    // Get the user token from the headers
    const auth_header = req.headers.authorization || '';

    // Extract username from auth header's information.
    var user_passwd = 'Anonymous:none';
    if (auth_header !== '') {
      user_passwd = Buffer.from(auth_header.split("Basic ")[1], 'base64').toString('utf-8');
    }
    var [user, password] = user_passwd.split(":"); // We don't care about the password in this example.

    // Offload authz decision to OPA.
    var url = opa_url + policy_path;
    var allowed = await opa_authz(
      url,
      req.body["query"],
      print(typeDefs),
      user,
      req.body["variables"],
      req.query["token"]);

    // If OPA rejects this request, we throw a generic error, and
    // log more detailed information for debugging.
    if (allowed) {
      console.log(`Success: user '${user}' is authorized \n`);
    } else {
      console.log(`Error: user '${user}' is not authorized to make query \n${req.body["query"]}\nwith variables:\n${JSON.stringify(req.body["variables"])}\n`);
      throw new AuthenticationError('Request rejected by server policy.');
    }

    // Add the username to the context.
    return { user };
  },
});


server.listen({ host: host, port: port }).then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});
