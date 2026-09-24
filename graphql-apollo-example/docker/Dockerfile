FROM node:18.4-buster-slim

WORKDIR /app
RUN npm install apollo-server graphql axios

COPY echo_server.js /app

ENV NODE_ENV=development
ENV HOST="0.0.0.0"
ENV PORT="5000"

CMD ["node", "echo_server.js"]

EXPOSE 5000
