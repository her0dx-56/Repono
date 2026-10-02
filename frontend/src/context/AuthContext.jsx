
import { createContext, useState } from 'react';
import { authService } from '../services/authService.js';

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => authService.getCurrent());
  const [loading, setLoading] = useState(false);

  async function login(data, mode = 'login') {
    setLoading(true);

    try {
      const authenticatedUser =
        mode === 'signup'
          ? await authService.signup(data)
          : await authService.login(data);

      setUser(authenticatedUser);

      return authenticatedUser;
    } finally {
      setLoading(false);
    }
  }

  async function logout() {
    await authService.logout();
    setUser(null);
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}