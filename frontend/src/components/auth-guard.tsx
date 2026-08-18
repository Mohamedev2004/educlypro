import { Navigate, Outlet, useLocation } from "react-router-dom"
import { useAuth } from "../context/auth/auth-context"
import {
  ONBOARDING_PATH,
  getPostLoginPath,
  needsOnboarding,
} from "@/utils/navigation-utils"

export function ProtectedRoute() {
  const { isAuthenticated, isLoading, user } = useAuth()
  const location = useLocation()

  if (isLoading) return null

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (needsOnboarding(user) && location.pathname !== ONBOARDING_PATH) {
    return <Navigate to={ONBOARDING_PATH} replace />
  }

  return <Outlet />
}

export function PublicOnlyRoute() {
  const { isAuthenticated, user } = useAuth()

  if (isAuthenticated) {
    return <Navigate to={getPostLoginPath(user)} replace />
  }

  return <Outlet />
}
