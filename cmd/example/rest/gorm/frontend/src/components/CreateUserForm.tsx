import React, { useState } from 'react'
import { createUser, UserDTO } from '../api'

type Props = {
  onCreated?: (u: UserDTO) => void
}

export default function CreateUserForm({ onCreated }: Props) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError(null)
    try {
      const u = await createUser({ name, email, password })
      setName('')
      setEmail('')
      setPassword('')
      if (onCreated) onCreated(u)
      // navigate to details by default
      window.location.hash = `#/users/${u.id}`
    } catch (err: any) {
      setError(err.message || 'failed to create')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="panel">
      <h3 className="panel-title">Add Person</h3>
      <form onSubmit={submit} className="form">
        <label className="field">
          <span className="field-label">Name</span>
          <input value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Full name" required />
        </label>

        <label className="field">
          <span className="field-label">Email</span>
          <input value={email} onChange={e => setEmail(e.target.value)} type="email" placeholder="email@example.org" required />
        </label>

        <label className="field">
          <span className="field-label">Password</span>
          <input value={password} onChange={e => setPassword(e.target.value)} type="password" placeholder="at least 8 characters" required minLength={8} />
        </label>

        {error && <div className="text-sm text-red-600">{error}</div>}

        <div className="form-actions" style={{ marginTop: 8 }}>
          <button type="submit" className="btn btn-primary" disabled={saving}>{saving ? 'Adding…' : 'Add Person'}</button>
          <button type="reset" onClick={() => { setName(''); setEmail(''); setPassword(''); }} className="btn btn-ghost">Reset</button>
        </div>
      </form>
    </div>
  )
}
