const fetch = require("sync-fetch");
const jsonGraphqlExpress = require("json-graphql-server").default;
const data = fetch(
  "https://github.com/rawkode/rawkode/releases/download/blox/data.json"
).json();

const serverless = require("serverless-http");
const app = require("express")();

const functionName = "serverless-http";

app.use("/", jsonGraphqlExpress(data));
console.log(data);

exports.handler = serverless(app);
