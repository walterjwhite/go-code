import React, { useState } from 'react'
import { updateUser, UserDTO } from '../api'

type Props = {
  user: UserDTO
  onSaved: (u: UserDTO) => void
  onCancel?: () => void
}

export default function UserForm({ user, onSaved, onCancel }: Props) {
  const [name, setName] = useState(user.name)
  const [email, setEmail] = useState(user.email)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError(null)
    try {
      const updated = await updateUser(user.id, { name, email })
      onSaved(updated)
    } catch (err: any) {
      setError(err.message || 'failed to save')
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="space-y-4">
      <div>
        <label className="block text-sm text-slate-600 mb-1">Name</label>
        <input
          value={name}
          onChange={e => setName(e.target.value)}
          className="w-full rounded-md border px-3 py-2 bg-white"
          required
          type="text"
          minLength={2}
        />
      </div>

      <div>
        <label className="block text-sm text-slate-600 mb-1">Email</label>
        <input
          value={email}
          onChange={e => setEmail(e.target.value)}
          className="w-full rounded-md border px-3 py-2 bg-white"
          type="email"
          required
        />
      </div>

      {error && <div className="text-sm text-red-600">{error}</div>}

      <div className="flex items-center gap-2 justify-end">
        <button
          type="button"
          onClick={onCancel}
          className="btn px-3 py-2 rounded-md bg-white border text-slate-700"
          disabled={saving}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="btn px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
          disabled={saving}
        >
          {saving ? 'Saving…' : 'Save'}
        </button>
      </div>
    </form>
  )
}