# The Next Tails Blog

This blog uses Nextjs + TailwindCSS for the front end and leverages Contentful as a headless CMS which you can use to write your blog.
Live Demo: https://next-tails-blog.vercel.app/

## How to use

### Update Env Variables

You can update your site metadata / social media in the next.config.js file.
You can also change othe variables such as how many blogs to show on the Recent Posts section
and update for any new routes for your sitemap.

```
env: {
  siteTitle: 'Next Blog',
  siteDescription: 'Next Tails Blog.',
  siteKeywords: 'nextjs, tailwindcss, contentful, blog',
  siteUrl: 'https://next-tails-blog.vercel.app/',
  siteImagePreviewUrl: '/images/main-img-preview.jpg',
  mainRoutes: ['/index', '/about', '/contact', '/blog'], // for sitemap; blog posts are generated dynamically
  blogRoute: '/blog', // for sitemap
  recentBlogNum: 3, // no. of blogs to display in recent posts
  twitterHandle: '@your_handle',
  twitterUrl: 'https://twitter.com',
  facebookUrl: 'https://facebook.com',
  instagramUrl: 'https://instagram.com',
  pinterestUrl: 'https://pinterest.com',
  youtubeUrl: 'https://youtube.com',
}
```

### Update Colors

You can update the color palette in tailwind.config.js file.

```
colors: {
  palette: {
    lighter: '',
    light: '',
    primary: '',
    dark: '',
  },
},
```

### Update Progressive Web App (PWA) data

Update the manifest.json file and the icons under the public/images/icons folder.
You can use free tools online such as https://realfavicongenerator.net/ to quickly generate all the different icon sizes and favicon.ico file.

<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/site.webmanifest">
<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#5bbad5">
<meta name="msapplication-TileColor" content="#da532c">
<meta name="theme-color" content="#ffffff">

### Setup Contentful

Setup a Contentful Account. After you setup your space, go to API keys under settings. Contentful should already have generated a couple of API Keys for you.
You will find the Space Id and Access Token under the Content Delivery API.
The url will look something like this: https://app.contentful.com/spaces/{YOUR_SPACE_ID}/api/keys.

Create a .env.local file on your project root and add the credentials there.

```
NEXT_PUBLIC_CONTENTFUL_SPACE_ID=''
NEXT_PUBLIC_CONTENTFUL_ACCESS_TOKEN=''
NEXT_PUBLIC_BLOG_COLLECTION=''
```

#### Your Blog Collection

This template uses one collection (or Content Type in Contentful) and works with the fields used in that collection.
You can follow the example and create similar fields at the start and then once you have everything up and running you can tinker
with the graphql queries in case you want to rename / add fields. You can name the collection whatever you want but be sure to add
it to the .env.local file as the NEXT_PUBLIC_BLOG_COLLECTION variable.

There are 8 fields:

```
Title
Slug
Hero Image
Description
Body
Author
Publish Date
Tags
```

![alt text](https://github.com/btahir/next-tailwind/blob/next-blog/public/images/contentful-collection.png)

Body is the only Rich Text field. Author and Tags are included but are not currently used on the site.

### Running Locally

Change into the project directory and run the following command:

```
yarn && yarn dev
```

### Deployment

You can deploy the site easily on Vercel and Netlify. Just hook up your Github with them.
Make sure you add the 3 environment variables in the .env.local file to your Vercel/Netlify Environment Variables when you deploy.

Once your site is up you can easily add Webhooks to Netlify/Vercel (Contentful has Templates for both so all you need to do is add in the hook url).
This will automatically redeploy the site whenever you publish/unpublish a new post and your site will always be in sync with your CMS.
