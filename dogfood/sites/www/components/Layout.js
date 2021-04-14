import Nav from '@/components/Nav';

function Layout({ children }) {
  return (
    <div className="max-w-7xl mx-auto sm:px-6 lg:px-8">
      <Nav> </Nav>
      {children}</div>
  )
}

export default Layout;
