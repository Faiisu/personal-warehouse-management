import type { User } from '../types/user'

type UserPageProps = {
  user?: User
  emailFallback?: string
}

const formatValue = (value?: string) => value || '—'
const initials = (value?: string) =>
  (value || '')
    .trim()
    .split(/\s+/)
    .map((part) => part[0]?.toUpperCase())
    .filter(Boolean)
    .slice(0, 2)
    .join('') || 'U'

function UserPage({ user, emailFallback }: UserPageProps) {
  const userId = user?.UserId ?? user?.UserID
  const displayName = user?.DisplayName || emailFallback || 'Your account'
  const avatarUrl = user?.AvatarURL
  const email = user?.Email ?? emailFallback
  const status = user?.Status

  return (
    <div className="main-card">
      <p className="badge">Account Detail</p>

      <div className="profile-hero">
        <div className="avatar" aria-label="profile avatar">
          {avatarUrl ? (
            <img src={avatarUrl} alt={displayName} />
          ) : (
            <span>{initials(displayName)}</span>
          )}
        </div>
        <div className="profile-meta">
          <h1 className="profile-name">{displayName}</h1>
          <p className='helper'>user-id: {userId}</p>         
          <div className="profile-pills">
            <span
              className={
                status && status.toUpperCase() === 'ACTIVE'
                  ? 'pill pill-active'
                  : 'pill subtle'
              }
            >
              Status: {formatValue(status)}
            </span>
            <span className="pill">Member</span>
          </div>
        </div>
      </div>

    </div>
  )
}

export default UserPage
