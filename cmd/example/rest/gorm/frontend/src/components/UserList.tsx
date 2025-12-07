import React, { useEffect, useMemo, useState } from 'react'
import { listUsers, getUser, deleteUser, UserDTO } from '../api'
import UserForm from './UserForm'
import CreateUserForm from './CreateUserForm'

export default function UserList() {
  const [users, setUsers] = useState<UserDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selected, setSelected] = useState<UserDTO | null>(null)
  const [query, setQuery] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)

  const pageSize = 20

  const fetchPage = async (p = 1) => {
    setLoading(true)
    setError(null)
    try {
      const data = await listUsers(p, pageSize)
      setUsers(data.items)
      setTotal(data.meta?.total ?? 0)
      setPage(p)
    } catch (err: any) {
      setError(err.message || 'failed to load')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchPage(1)
  }, [])

  const onEdit = async (id: number) => {
    try {
      const u = await getUser(id)
      setSelected(u)
    } catch (err: any) {
      setError(err.message || 'failed to fetch user')
    }
  }

  const onView = (id: number) => {
    // navigate to details page (hash-based)
    window.location.hash = `#/users/${id}`
  }

  const onDelete = async (id: number) => {
    if (!confirm('Delete this user? This action cannot be undone.')) return
    try {
      await deleteUser(id)
      // optimistic refresh
      setUsers(v => v.filter(u => u.id !== id))
      setTotal(t => t - 1)
    } catch (err: any) {
      setError(err.message || 'failed to delete')
    }
  }

  const onSaved = (updated: UserDTO) => {
    // update in list (if present)
    setUsers(prev => prev.map(u => (u.id === updated.id ? updated : u)))
    setSelected(null)
  }

  const pageCount = useMemo(() => Math.max(1, Math.ceil(total / pageSize)), [total])

  const filteredUsers = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return users
    return users.filter(u => u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q))
  }, [users, query])

  return (
    <div className="panel">
      <div className="flex items-center justify-between mb-4 gap-12">
        <div>
          <h2 className="text-lg font-medium m-0">Users</h2>
          <div className="text-sm text-slate-500">Total: {total}</div>
        </div>
        <div className="flex gap-8 items-center">
          <button onClick={() => alert('Create flow: use API or request UI addition')} className="btn btn-ghost">Create</button>
        </div>
      </div>

      <div className="mb-12">
        <label className="block text-sm text-slate-600 mb-1">Search (name or email)</label>
        <input type="text" value={query} onChange={e => setQuery(e.target.value)} placeholder="Search by name or email" className="w-full rounded-md border px-3 py-2 bg-white" />
      </div>
      <div className="flex-gap-18">
        <div className="flex-1-1-520">
          <div className="overflow-x-auto">
            <table className="w-full text-left table-auto">
              <thead className="bg-slate-50">
                <tr>
                  <th className="px-4 py-2 text-sm text-slate-500">Name</th>
                  <th className="px-4 py-2 text-sm text-slate-500">Email</th>
                  <th className="px-4 py-2 text-sm text-slate-500 w-36">Actions</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr>
                    <td colSpan={3} className="px-4 py-6 text-center text-slate-400">Loading...</td>
                  </tr>
                ) : users.length === 0 ? (
                  <tr>
                    <td colSpan={3} className="px-4 py-6 text-center text-slate-500">No users found.</td>
                  </tr>
                ) : (
                  filteredUsers.map(u => (
                      <tr key={u.id} className={`border-t transition-colors odd:bg-white even:bg-slate-50 hover:bg-slate-100`}>
                        <td className="px-4 py-3">
                          <a href={`#/users/${u.id}`} className="block font-medium text-left no-underline text-slate-900">
                            {u.name}
                          </a>
                        </td>
                        <td className="px-4 py-3 text-slate-600">
                          <a href={`#/users/${u.id}`} className="block no-underline text-slate-600">
                            {u.email}
                          </a>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex gap-2">
                            <button onClick={() => onDelete(u.id)} className="btn text-sm px-3 py-1 bg-red-50 border border-red-200 text-red-700 rounded-md hover:bg-red-100">Delete</button>
                          </div>
                        </td>
                      </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          <div className="mt-4 flex items-center justify-between">
            <div className="text-sm text-slate-500">Page {page} of {pageCount}</div>
            <div className="flex gap-2">
              <button onClick={() => fetchPage(Math.max(1, page - 1))} disabled={page <= 1} className="btn px-3 py-1 rounded-md bg-white border text-slate-700 disabled:opacity-50">Prev</button>
              <button onClick={() => fetchPage(Math.min(pageCount, page + 1))} disabled={page >= pageCount} className="btn px-3 py-1 rounded-md bg-white border text-slate-700 disabled:opacity-50">Next</button>
            </div>
          </div>
        </div>

        <aside className="flex-0-0-320">
          <CreateUserForm onCreated={(u) => { setUsers(prev => [u, ...prev]); setTotal(t => t + 1) }} />
          <div className="panel mt-12">
            <h4 className="mb-2">Tips</h4>
            <ul style={{ margin: 0, paddingLeft: 18 }}>
              <li>Click a name to view details.</li>
              <li>Use the Delete button to remove an entry.</li>
            </ul>
          </div>
        </aside>
      </div>

      {/* Drawer/modal area */}
      {selected && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setSelected(null)} />
          <div className="relative bg-white rounded-t-xl sm:rounded-xl shadow-lg w-full sm:max-w-xl mx-4 my-6 p-6 z-10">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-medium">Edit User</h3>
              <button className="text-slate-400" onClick={() => setSelected(null)}>✕</button>
            </div>
            <UserForm user={selected} onSaved={onSaved} onCancel={() => setSelected(null)} />
          </div>
        </div>
      )}

      {error && (
        <div className="mt-4 text-sm text-red-600">{error}</div>
      )}
    </div>
  )
}