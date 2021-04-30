
# Serving Data over REST

## Prerequisites

* content repository hosted on GitHub
* blox.cue configuration is complete
* [dataset hosted on GitHub](/recipes/github-releases), or another fixed URL

## json-server

Use the awesome [json-server](https://www.npmjs.com/package/json-server) to serve your dataset as a REST API.

The implementation will vary based on how you choose to run the service, but we've really enjoyed using `serverless`/Function hosting platforms like Azure Functions, Vercel, AWS Lambda, and Netlify to host this service.  Here's the core of the recipe:

```javascript
const fetch = require("sync-fetch");

const jsonServer = require('json-server')
const express = require("express");

const data = fetch(
  "https://github.com/you/yourrepo/releases/download/blox/data.json"
).json();
const app = require("express")();
const router = jsonServer.router(data, { foreignKeySuffix: '_id' })


app.use("/api", router);


const port = process.env.PORT || 3000;

module.exports = app.listen(port, () =>
  console.log(`Server running on ${port}, http://localhost:${port}`)
);
```

We're using `json-server` to serve the dataset as an `express` application at the `/api` route. When you run this, you can make a `GET` request to `/api/datasetname` (`/api/articles` for example) and get back all the data in that dataset.

See the [documentation](https://www.npmjs.com/package/json-server) for complete details.

This recipe won't typically work without some modification of the source. The key is to figure out how to adapt an `express` endpoint to your hosting provider's serverless hosting. Vercel is happy to serve up a function running `express` without modification. Azure Functions requires an adapter like [azure-function-express](https://www.npmjs.com/package/azure-function-express).


