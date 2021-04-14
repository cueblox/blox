import { ApolloClient, InMemoryCache, gql } from '@apollo/client';

const client = new ApolloClient({
  uri: process.env.NEXT_PUBLIC_CMS_URL || 'https://api.cueblox.com/api/graphql',
  cache: new InMemoryCache()
});



export async function getPage(slug) {

  const { data } = await client.query({

    query: gql`
{
  allPages(filter: {id:"${slug}"}) {
    id
    title
    
  }
}
`});

  const pages = data.allPages;
  console.log(pages)
  if (pages.length > 0) {
    return pages[0]
  }

}


export async function getSections() {

  const { data } = await client.query({

    query: gql`
{
  allSections(sortField: "weight", sortOrder: "asc") {
    id
    name
    description
    weight
    Pages{
      id
      title
      weight
    }
  }
}
`});

  const sections = data.allSections;

  return sections


}

export async function getArticles() {
  const { data } = await client.query({
    query: gql`
    query {
    allArticles {
    id
        title
        Category{
    id
          name
  }
        tags
        publish_date
        last_edit_date
        edit_description
      }
    }
`
  });
  const articles = data.allArticles;
  return articles;
}
