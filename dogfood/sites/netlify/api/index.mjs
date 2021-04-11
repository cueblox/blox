import express from 'express';
import fetch from 'sync-fetch';
import jsonGraphqlExpress from 'json-graphql-server';
import serverless from 'serverless-http';

const data = fetch(
  "https://github.com/cueblox/blox/releases/download/blox/data.json"
).json();

const app = express();

const functionName = "serverless-http";

app.use("/", jsonGraphqlExpress(data));
console.log(data);

exports.handler = serverless(app);
