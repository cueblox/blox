const fetch = require('sync-fetch');

const jsonGraphqlExpress = require('json-graphql-server').default;
const createHandler = require("azure-function-express").createHandler;
const jsonServer = require('json-server')


const data = fetch('https://github.com/cueblox/blox/releases/download/blox/data.json').json();

const router = jsonServer.router(data, { foreignKeySuffix: '_id' })
const app = require('express')();


app.use('/api/graphql', jsonGraphqlExpress(data));
app.use("/api", router);

console.log(data);

module.exports = createHandler(app);
