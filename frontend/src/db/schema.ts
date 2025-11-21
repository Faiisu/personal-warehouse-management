// Shared user collection shape for backend usage and front-end typing.
export const usersCollection = 'users'

export type UserStatus = 'ACTIVE' | 'INACTIVE'

export interface UserDoc {
  _id?: string
  email: string
  display_name: string
  password_hash: string
  avatar_url?: string | null
  status: UserStatus
}

// Suggested indexes to mirror your schema intent.
export const usersIndexes = [
  { key: { email: 1 }, unique: true },
]
