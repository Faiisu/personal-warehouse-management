import { useState } from 'react'
import './App.css'
import AuthFormCard from './components/AuthForm'
import MainPage from './components/MainPage'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [userEmail, setUserEmail] = useState<string | undefined>(undefined)

  const handleLoginSuccess = (
    _resp?: unknown,
    context?: { email: string },
  ) => {
    setIsAuthenticated(true)
    setUserEmail(context?.email)
  }

  const handleLogout = () => {
    setIsAuthenticated(false)
    setUserEmail(undefined)
  }

  return (
    <div className="page">
      {isAuthenticated ? (
        <MainPage email={userEmail} onLogout={handleLogout} />
      ) : (
        <AuthFormCard onLoginSuccess={handleLoginSuccess} />
      )}
    </div>
  )
}

export default App
