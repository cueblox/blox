import Nav from './Nav'

/* This example requires Tailwind CSS v2.0+ */


function classNames(...classes) {
  return classes.filter(Boolean).join(' ')
}

export default function Layout({ children }) {

  return (
    <div className="h-screen flex overflow-hidden bg-gray-100">
      <Nav></Nav>
      <main className="flex-1 relative z-0 overflow-y-auto focus:outline-none">
        <div className="py-6">

          <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
            {children}
            {/* /End replace */}
          </div>
        </div>
      </main>
    </div>
  )
}
