export type UserDTO = {
  id: number
  name: string
  email: string
}

const BASE = (import.meta.env.VITE_API_URL as string) || ''

function url(path: string) {
  return BASE + path
}

export async function listUsers(page = 1, size = 20) {
  const res = await fetch(url(`/api/users?page=${page}&size=${size}`))
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<{ items: UserDTO[]; meta: { total: number } }>
}

export async function getUser(id: number) {
  const res = await fetch(url(`/api/users/${id}`))
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<UserDTO>
}

export async function updateUser(id: number, data: Partial<{ name: string; email: string }>) {
  const res = await fetch(url(`/api/users/${id}`), {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<UserDTO>
}

export async function deleteUser(id: number) {
  const res = await fetch(url(`/api/users/${id}`), { method: 'DELETE' })
  if (res.status === 204) return
  if (!res.ok) throw new Error(await res.text())
}

export async function createUser(data: { name: string; email: string; password: string }) {
  const res = await fetch(url(`/api/users`), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<UserDTO>
}