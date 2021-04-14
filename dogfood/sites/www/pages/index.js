import { GlobalContext } from "./_app";
import { getPage } from "@/lib/cms";

export default function Index() {
  return (
    <main className="mt-16 mx-auto max-w-7xl px-4 sm:mt-24">
      <div className="text-center">
        <h1 className="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
          <span className="block xl:inline">Do more with your</span>{' '}
          <span className="block text-royalblue-600 xl:inline">YAML and Markdown</span>
        </h1>
        <p className="mt-3 max-w-md mx-auto text-base text-gray-500 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
          Coming Soon.
                </p>
        <div className="mt-5 max-w-md mx-auto sm:flex sm:justify-center md:mt-8">
          <div className="rounded-md shadow">
            <a
              href="/docs"
              className="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-denim-600 hover:bg-denim-700 md:py-4 md:text-lg md:px-10"
            >
              Get started
                    </a>
          </div>
          <div className="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
            <a
              href="/youtube"
              className="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-denim-600 bg-white hover:bg-gray-50 md:py-4 md:text-lg md:px-10"
            >
              Demo Video
                    </a>
          </div>
        </div>
      </div>
    </main>

  )
}

export async function getStaticProps() {

  const page = await getPage("index");

  console.log(page)

  return {
    props: {
      page

    },
  }
}
function merge(a, b, prop) {
  var reduced = a.filter(aitem => !b.find(bitem => aitem[prop] === bitem[prop]))
  return reduced.concat(b);
}