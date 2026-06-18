import React, { useState, useEffect } from 'react'
import Login from './components/Login.jsx'
import Lobby from './components/Lobby.jsx'
import PokerTable from './components/PokerTable.jsx'
import { LogOut } from 'lucide-react'

const API_URL = `http://${window.location.hostname}:8081`;

export default function App() {
  const [token, setToken] = useState(null)
  const [user, setUser] = useState(null)
  const [activeTableId, setActiveTableId] = useState(null)
  const [appLoading, setAppLoading] = useState(true)

  // Load auth state from localStorage on init
  useEffect(() => {
    const savedToken = localStorage.getItem('auth_token')
    const savedUser = localStorage.getItem('auth_user')

    if (savedToken && savedUser) {
      setToken(savedToken)
      setUser(JSON.parse(savedUser))
    }
    setAppLoading(false)
  }, [])

  const handleAuthSuccess = (newToken, newUser) => {
    setToken(newToken)
    setUser(newUser)
  }

  const handleLogout = async () => {
    if (token) {
      try {
        await fetch(`${API_URL}/api/logout`, {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` }
        })
      } catch (err) {
        console.error("Erreur lors de la déconnexion :", err)
      }
    }
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
    setToken(null)
    setUser(null)
    setActiveTableId(null)
  }

  // Refresh user money and profile info
  const refreshUserProfile = async () => {
    if (!token) return
    try {
      const response = await fetch(`${API_URL}/api/users/me`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (response.ok) {
        const freshUser = await response.json()
        setUser(freshUser)
        localStorage.setItem('auth_user', JSON.stringify(freshUser))
      }
    } catch (err) {
      console.error("Impossible d'actualiser le profil :", err)
    }
  }

  const handleJoinTable = (tableId) => {
    setActiveTableId(tableId)
  }

  const handleLeaveTable = () => {
    setActiveTableId(null)
    refreshUserProfile() // Refresh chips balance on cash-out return
  }

  if (appLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', background: 'var(--bg-dark)' }}>
        <h2 style={{ color: 'var(--accent-gold)' }}>Chargement de l'application...</h2>
      </div>
    )
  }

  return (
    <div className="app-container">
      {/* Header bar */}
      <header className="main-header">
        <div className="brand">
          <span className="brand-icon">♣️</span>
          <span>Poker Fil Violet</span>
        </div>

        {token && user && (
          <div className="user-profile-widget">
            <div className="chip-badge">
              <span className="chip-icon"></span>
              <span>{user.money} Jetons</span>
            </div>
            
            <span style={{ fontWeight: '600', fontSize: '0.95rem' }}>
              {user.nick_name}
            </span>

            <button className="btn btn-secondary" onClick={handleLogout} title="Se Déconnecter">
              <LogOut size={16} />
              <span style={{ fontSize: '0.85rem' }}>Déconnexion</span>
            </button>
          </div>
        )}
      </header>

      {/* View routing */}
      <main style={{ flex: 1 }}>
        {!token ? (
          <Login onAuthSuccess={handleAuthSuccess} />
        ) : activeTableId ? (
          <PokerTable 
            token={token} 
            user={user} 
            tableId={activeTableId} 
            onLeaveTable={handleLeaveTable} 
          />
        ) : (
          <Lobby 
            token={token} 
            user={user} 
            onJoinTable={handleJoinTable} 
          />
        )}
      </main>
    </div>
  )
}
