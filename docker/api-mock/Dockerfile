FROM node:12-alpine

ARG projectName=api-mock

RUN echo ${projectName}

WORKDIR /app

RUN npm install -g json-server

EXPOSE 3000
# Example
# CMD [ "json-server","--host", "0.0.0.0","--watch","/app/samples/api-mock/example.db.json", "--routes","/app/samples/api-mock/example.routes.json","--id","sku" ]
CMD [ "json-server","--host", "0.0.0.0","--watch","/app/samples/api-mock/sefaz.db.json", "--routes","/app/samples/api-mock/sefaz.routes.json" ]