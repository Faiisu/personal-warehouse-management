import { useEffect, useState } from 'react'
import './App.css'
import AuthFormCard from './components/AuthForm'
import UserPage from './components/UserPage'
import AdminDashboard from './components/AdminDashboard'
import StockPage from './components/StockPage'
import StockProductsPage from './components/StockProductsPage'
import type { User } from './types/user'

const LOCAL_STORAGE_USER_KEY = 'authUser'

type ViewKey = 'events' | 'create' | 'stocks' | 'user'

function App() {
  const pathname =
    typeof window !== 'undefined' && window.location.pathname
      ? window.location.pathname
      : '/'
  const normalizedPath = pathname.toLowerCase()
  const isAdminRoute =
    typeof window !== 'undefined' && normalizedPath.startsWith('/admin')
  const isAccountRoute =
    typeof window !== 'undefined' &&
    (normalizedPath === '/account' || normalizedPath === '/account/')
  const isStockDetailRoute =
    !isAdminRoute &&
    !isAccountRoute &&
    typeof window !== 'undefined' &&
    normalizedPath.startsWith('/stocks/') &&
    normalizedPath.length > '/stocks/'.length
  const stockSlug = isStockDetailRoute
    ? decodeURIComponent(pathname.replace(/^\/stocks\//i, ''))
    : ''

  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [userEmail, setUserEmail] = useState<string | undefined>(undefined)
  const [user, setUser] = useState<User | undefined>(undefined)
  const [activeView, setActiveView] = useState<ViewKey>('stocks')
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    const cached = localStorage.getItem(LOCAL_STORAGE_USER_KEY)
    if (!cached) {
      setHydrated(true)
      return
    }
    try {
      const parsed = JSON.parse(cached) as User
      setUser(parsed)
      setUserEmail(parsed.Email)
      setIsAuthenticated(true)
    } catch {
      localStorage.removeItem(LOCAL_STORAGE_USER_KEY)
    } finally {
      setHydrated(true)
    }
    return () => {
      setHydrated(true)
    }
  }, [])

  useEffect(() => {
    if (
      typeof window !== 'undefined' &&
      !isAdminRoute &&
      !isStockDetailRoute &&
      window.location.pathname === '/'
    ) {
      window.location.replace('/stocks')
    }
    if (isAccountRoute) {
      setActiveView('user')
    } else {
      setActiveView('stocks')
    }
    if (
      typeof window !== 'undefined' &&
      !isAdminRoute &&
      !isAccountRoute &&
      !isStockDetailRoute &&
      normalizedPath !== '/' &&
      normalizedPath !== '/stocks' &&
      normalizedPath !== '/stocks/' &&
      normalizedPath !== '/login' &&
      normalizedPath !== '/login/' &&
      !normalizedPath.startsWith('/stocks/')
    ) {
      const suffix = pathname.startsWith('/') ? pathname : `/${pathname}`
      window.location.replace(`/stocks${suffix}`)
    }
  }, [
    isAdminRoute,
    isStockDetailRoute,
    isAccountRoute,
    normalizedPath,
    pathname,
  ])

  useEffect(() => {
    if (
      isAuthenticated &&
      hydrated &&
      typeof window !== 'undefined' &&
      normalizedPath === '/login'
    ) {
      window.location.replace('/stocks')
    }
  }, [isAuthenticated, normalizedPath, hydrated])

  const handleLoginSuccess = (
    loggedInUser?: User,
    context?: { email: string },
  ) => {
    setIsAuthenticated(true)
    setUser(loggedInUser)
    const email = loggedInUser?.Email ?? context?.email
    setUserEmail(email)

    if (loggedInUser) {
      localStorage.setItem(LOCAL_STORAGE_USER_KEY, JSON.stringify(loggedInUser))
    } else if (email) {
      localStorage.setItem(
        LOCAL_STORAGE_USER_KEY,
        JSON.stringify({ Email: email }),
      )
    }
    if (typeof window !== 'undefined') {
      window.location.replace('/stocks')
    }
  }

  const handleLogout = () => {
    setIsAuthenticated(false)
    setUserEmail(undefined)
    setUser(undefined)
    localStorage.removeItem(LOCAL_STORAGE_USER_KEY)
    setActiveView('events')
  }

  if (!isAuthenticated) {
    if (
      hydrated &&
      typeof window !== 'undefined' &&
      normalizedPath !== '/login'
    ) {
      window.location.replace('/login')
      return null
    }
    return (
      <div className="page centered">
        <AuthFormCard onLoginSuccess={handleLoginSuccess} />
        {isAdminRoute && (
          <p className="subhead">Admin access requires login.</p>
        )}
      </div>
    )
  }

  const viewLabel: Record<ViewKey, string> = {
    stocks: 'Stocks',
    user: 'Account',
  }

  let mainContent: JSX.Element
  if (isAdminRoute) {
    mainContent = (
      <AdminDashboard user={user} email={userEmail} onLogout={handleLogout} />
    )
  } else if (isStockDetailRoute) {
    mainContent = (
      <StockProductsPage
        stockName={stockSlug}
        user={user}
        onBack={() => window.location.assign('/stocks')}
      />
    )
  } else if (isAccountRoute || activeView === 'user') {
    mainContent = <UserPage user={user} emailFallback={userEmail} />
  } else {
    mainContent = <StockPage user={user} />
  }
  const handleNav = (view: ViewKey) => {
    if (view === 'user') {
      window.location.assign('/account')
    } else if (view === 'stocks') {
      window.location.assign('/stocks')
    }
    setActiveView(view)
  }

  return (
    <div className="page shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-dot" />
          <div>
            <p className="brand-name">{user?.DisplayName}</p>
          </div>
        </div>
        <nav className="nav">
          {(Object.keys(viewLabel) as ViewKey[]).map((key) => (
            <button
              key={key}
              type="button"
              className={activeView === key ? 'nav-item active' : 'nav-item'}
              onClick={() => handleNav(key)}
            >
              {viewLabel[key]}
            </button>
          ))}
        </nav>
        <div className="nav-divider" />
        <button type="button" className="outline full" onClick={handleLogout}>
          Log out
        </button>
      </aside>

      <div className="main-area">
        
        <div className="mobile-nav">
          <select
            value={activeView}
            onChange={(e) => handleNav(e.target.value as ViewKey)}
          >
            {(Object.keys(viewLabel) as ViewKey[]).map((key) => (
              <option key={key} value={key}>
                {viewLabel[key]}
              </option>
            ))}
          </select>
          {/* <button type="button" className="outline" onClick={handleLogout}>
            Log out
          </button> */}
        </div>

        <main className="content">{mainContent}</main>
      </div>
    </div>
  )
}

export default App
