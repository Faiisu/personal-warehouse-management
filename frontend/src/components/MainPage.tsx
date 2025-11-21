type MainPageProps = {
  email?: string
  onLogout?: () => void
}

function MainPage({ email, onLogout }: MainPageProps) {
  return (
    <div className="main-card">
      <p className="badge">Event Blog App</p>
      <h1>Welcome back{email ? `, ${email}` : ''}!</h1>
      <p className="subhead">
        You are now logged in. This is your main page to manage events and posts.
      </p>

      <div className="main-actions">
        <button type="button" className="submit">
          Go to dashboard
        </button>
        <button
          type="button"
          className="outline"
          onClick={onLogout}
        >
          Log out
        </button>
      </div>
    </div>
  )
}

export default MainPage
