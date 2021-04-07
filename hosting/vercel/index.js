const fetch = require('sync-fetch')

import jsonGraphqlExpress from 'json-graphql-server';

const data = fetch('https://github.com/rawkode/rawkode/releases/download/blox/data.json').json();
const app = require('express')();


app.use('/api/graphql', jsonGraphqlExpress(data));
console.log(data);

const port = process.env.PORT || 3000;

module.exports = app.listen(port, () => console.log(`Server running on ${port}, http://localhost:${port}`));