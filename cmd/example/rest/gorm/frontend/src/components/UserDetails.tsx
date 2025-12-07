import React, { useEffect, useState } from 'react'
import { getUser, deleteUser, UserDTO } from '../api'
import UserForm from './UserForm'

type Props = { id: number }

export default function UserDetails({ id }: Props) {
  const [user, setUser] = useState<UserDTO | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [editing, setEditing] = useState(false)

  const fetchUser = async () => {
    setLoading(true)
    setError(null)
    try {
      const u = await getUser(id)
      setUser(u)
    } catch (err: any) {
      setError(err.message || 'failed to load')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUser()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const onDelete = async () => {
    if (!confirm('Delete this user? This action cannot be undone.')) return
    try {
      await deleteUser(id)
      // go back to list
      window.location.hash = '#/'
    } catch (err: any) {
      setError(err.message || 'failed to delete')
    }
  }

  const onSaved = (u: UserDTO) => {
    setUser(u)
    setEditing(false)
  }

  if (loading) return <div>Loading…</div>
  if (error && !user) return <div className="text-sm text-red-600">{error}</div>
  if (!user) return <div className="text-sm text-slate-600">User not found.</div>

  return (
    <div className="panel">
      <div className="flex items-start justify-between mb-4 gap-12">
        <div>
          <h2 className="text-xl font-semibold m-0">{user.name}</h2>
          <div className="text-sm text-slate-500">ID: {user.id}</div>
        </div>

        <div className="flex gap-8 items-center">
          <button onClick={() => fetchUser()} className="btn btn-outline btn-small">Refresh</button>
          <button onClick={() => (window.location.hash = '#/')} className="btn btn-ghost btn-small">Directory</button>
          <button onClick={onDelete} className="btn btn-danger btn-small">Delete</button>
        </div>
      </div>

      <hr className="hr-sep" />

      {error && <div className="alert alert-error mb-3" role="alert">{error}</div>}

      <div className="flex-gap-18">
        <div className="flex-1-1-320">
          <div className="panel">
            <h3 className="mb-2">Edit Contact</h3>
            {editing ? (
              <UserForm user={user} onSaved={onSaved} onCancel={() => setEditing(false)} />
            ) : (
              <div className="space-y-2">
                <div>
                  <div className="text-sm text-slate-500">Name</div>
                  <div className="font-medium">{user.name}</div>
                </div>
                <div>
                  <div className="text-sm text-slate-500">Email</div>
                  <div className="font-medium">{user.email}</div>
                </div>
                <div className="mt-2">
                  <button onClick={() => setEditing(true)} className="btn btn-primary">Edit</button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
