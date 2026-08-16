/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { useNavigate } from "react-router-dom"
import { useQueryClient } from "@tanstack/react-query"
import { AuthService } from "../../api/services/auth-service"
import { session } from "@/utils/session-utils"
import { useMe } from "../../hooks/use-me"
import type {
  LoginPayload,
  UpdateProfilePayload,
  UpdatePasswordPayload,
  User,
} from "../../api/types/auth.types"

type AuthContextValue = {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (payload: LoginPayload) => Promise<User>
  updateProfile: (payload: UpdateProfilePayload) => Promise<User>
  updatePassword: (payload: UpdatePasswordPayload) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = React.createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { data: user = null, isLoading } = useMe()

  // Handles force-logout when refresh token expires (fired from axios interceptor)
  React.useEffect(() => {
    const handleForceLogout = () => {
      queryClient.setQueryData(["me"], null)
      queryClient.clear()
      navigate("/login", { replace: true })
    }

    window.addEventListener("auth:logout", handleForceLogout)
    return () => window.removeEventListener("auth:logout", handleForceLogout)
  }, [navigate, queryClient])

  const login = React.useCallback(async (payload: LoginPayload) => {
    const currentUser = await AuthService.login(payload)
    queryClient.setQueryData(["me"], currentUser)
    return currentUser
  }, [queryClient])

  const updateProfile = React.useCallback(async (payload: UpdateProfilePayload) => {
    const updatedUser = await AuthService.updateProfile(payload)
    queryClient.setQueryData(["me"], updatedUser)
    return updatedUser
  }, [queryClient])

  const updatePassword = React.useCallback(async (payload: UpdatePasswordPayload) => {
    await AuthService.updatePassword(payload)
  }, [])

  const logout = React.useCallback(async () => {
    try {
      await AuthService.logout()
    } finally {
      session.clear()
      queryClient.setQueryData(["me"], null)
      queryClient.clear()
      navigate("/login", { replace: true })
    }
  }, [navigate, queryClient])

  const value = React.useMemo(() => ({
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    updateProfile,
    updatePassword,
    logout,
  }), [user, isLoading, login, updateProfile, updatePassword, logout])

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = React.useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used inside AuthProvider")
  return context
}
