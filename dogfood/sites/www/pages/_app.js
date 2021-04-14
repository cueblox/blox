import '@/styles/globals.css';

import App from "next/app";
import Layout from '@/components/Layout';
import { createContext } from "react";
import { getSections } from "@/lib/cms";

export const GlobalContext = createContext({});


function MyApp({ Component, pageProps }) {
  const { sections } = pageProps;
  return (
    <GlobalContext.Provider value={sections}>
      <Layout>
        <Component {...pageProps} />
      </Layout>
    </GlobalContext.Provider>

  )
}


// getInitialProps disables automatic static optimization for pages that don't
// have getStaticProps. So article, category and home pages still get SSG.
// Hopefully we can replace this with getStaticProps once this issue is fixed:
// https://github.com/vercel/next.js/discussions/10949
MyApp.getInitialProps = async (ctx) => {
  // Calls page's `getInitialProps` and fills `appProps.pageProps`
  const appProps = await App.getInitialProps(ctx);
  // Fetch global site settings from Strapi
  const sections = await getSections();
  // Pass the data to our page via props
  return { ...appProps, pageProps: { sections } };
};

export default MyApp;