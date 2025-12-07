import React, { useEffect, useState } from 'react'
import UserList from './components/UserList'
import UserDetails from './components/UserDetails'

// Simple hash-based routing used to avoid adding a router dependency.
// Routes:
//  - "" or "#/" -> list
//  - "#/users/:id" -> details

function parseHash(): { name: 'list' } | { name: 'details'; id: number } {
  const h = (window.location.hash || '').replace(/^#/, '') || '/'
  const parts = h.split('/').filter(Boolean)
  if (parts[0] === 'users' && parts[1]) {
    const id = Number(parts[1])
    if (!Number.isNaN(id)) return { name: 'details', id }
  }
  return { name: 'list' }
}

export default function App() {
  const [route, setRoute] = useState(parseHash())

  useEffect(() => {
    const onHash = () => setRoute(parseHash())
    window.addEventListener('hashchange', onHash)
    return () => window.removeEventListener('hashchange', onHash)
  }, [])

  return (
    <div className="min-h-screen text-slate-900">
      <header className="site-header">
        <div className="container">
          <a className="brand">Contacts — Users</a>
          <div className="site-nav">
            <div className="header-tagline">Simple CRUD UI</div>
          </div>
        </div>
      </header>

      <main className="container layout">
        <section className="list-panel">
          {route.name === 'list' ? <UserList /> : <UserDetails id={route.id} />}
        </section>

        <aside className="form-panel">
          {/* Right column is handled inside UserList (create panel) or UserDetails shows actions */}
        </aside>
      </main>

      <footer className="footer">
        <div className="container">
          <div className="small">Built with a modern React + Tailwind stack — connects to /api/users</div>
        </div>
      </footer>
    </div>
  )
}