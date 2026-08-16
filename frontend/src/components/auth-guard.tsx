import { Navigate, Outlet } from "react-router-dom"
import { useAuth } from "../context/auth/auth-context"
import { getDashboardPath } from "@/utils/navigation-utils"
import type { User } from "@/api/types/auth.types"

function getRedirectPath(user: User | null) {
  if (!user) return "/login"
  return getDashboardPath(user.role)
}

export function ProtectedRoute() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) return null

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}

export function PublicOnlyRoute() {
  const { isAuthenticated, user } = useAuth()

  if (isAuthenticated) {
    return (
      <Navigate
        to={getRedirectPath(user)}
        replace
      />
    )
  }

  return <Outlet />
}
