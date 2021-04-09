const fetch = require('sync-fetch');

const jsonGraphqlExpress = require('json-graphql-server').default;
const createHandler = require("azure-function-express").createHandler;

const data = fetch('https://github.com/cueblox/blox/releases/download/blox/data.json').json();
const app = require('express')();


app.use('/api/graphql', jsonGraphqlExpress(data));
console.log(data);

module.exports = createHandler(app);
