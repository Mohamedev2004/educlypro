import { Navigate, Route, Routes } from "react-router-dom"
import { lazy, Suspense } from "react"

import { ProtectedRoute, PublicOnlyRoute } from "./components/auth-guard"

// ======================
// Lazy Loaded Pages
// ======================

// Auth
const Login = lazy(() => import("./pages/auth/login"))
const ForgotPassword = lazy(() => import("./pages/auth/forgot-password"))
const ResetPassword = lazy(() => import("./pages/auth/reset-password"))

// Settings
const Profile = lazy(() => import("./pages/settings/profile"))
const Password = lazy(() => import("./pages/settings/password"))
const Appearance = lazy(() => import("./pages/settings/appearance"))

// Shared pages (reused across roles)
const Dashboard = lazy(() => import("./pages/dashboard"))
const NotificationsPage = lazy(() => import("./pages/notifications"))

// Super Admin
const SuperAdminLogs = lazy(() => import("./pages/super-admin/logs"))
const SuperAdminCenters = lazy(() => import("./pages/super-admin/centers"))
const SuperAdminCenterDetail = lazy(
  () => import("./pages/super-admin/center-detail")
)

// Center Owner
const CenterOwnerOnboarding = lazy(
  () => import("./pages/center-owner/onboarding")
)
const CenterOwnerStaff = lazy(() => import("./pages/center-owner/staff"))
const CenterOwnerSubCenters = lazy(
  () => import("./pages/center-owner/subcenters")
)
const CenterOwnerTeachers = lazy(() => import("./pages/center-owner/teachers"))
const CenterOwnerTeacherDetail = lazy(
  () => import("./pages/center-owner/teacher-detail")
)

// Center Scanner
const CenterScannerScanner = lazy(
  () => import("./pages/center-scanner/scanner")
)
const CenterScannerStudents = lazy(
  () => import("./pages/center-scanner/students")
)

export function App() {
  return (
    <Suspense fallback={null}>
      <Routes>
        {/* Public routes */}
        <Route element={<PublicOnlyRoute />}>
          <Route path="/" element={<Navigate to="/login" replace />} />
          <Route path="/login" element={<Login />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />
          <Route path="/reset-password" element={<ResetPassword />} />
        </Route>

        {/* Protected routes */}
        <Route element={<ProtectedRoute />}>
          {/* Settings */}
          <Route path="/profile" element={<Profile />} />
          <Route path="/password" element={<Password />} />
          <Route path="/appearance" element={<Appearance />} />

          {/* Super Admin */}
          <Route
            path="/super-admin/dashboard"
            element={<Dashboard dashboardHref="/super-admin/dashboard" />}
          />
          <Route
            path="/super-admin/notifications"
            element={
              <NotificationsPage
                dashboardHref="/super-admin/dashboard"
                breadcrumbTaglineKey="roles.overview"
                breadcrumbLabelKey="roles.notifications"
              />
            }
          />
          <Route path="/super-admin/logs" element={<SuperAdminLogs />} />
          <Route path="/super-admin/centers" element={<SuperAdminCenters />} />
          <Route
            path="/super-admin/centers/:slug"
            element={<SuperAdminCenterDetail />}
          />

          {/* Center Owner */}
          <Route
            path="/center-owner/onboarding"
            element={<CenterOwnerOnboarding />}
          />
          <Route
            path="/center-owner/dashboard"
            element={<Dashboard dashboardHref="/center-owner/dashboard" />}
          />
          <Route
            path="/center-owner/notifications"
            element={
              <NotificationsPage
                dashboardHref="/center-owner/dashboard"
                breadcrumbTaglineKey="roles.overview"
                breadcrumbLabelKey="roles.notifications"
              />
            }
          />
          <Route path="/center-owner/staff" element={<CenterOwnerStaff />} />
          <Route
            path="/center-owner/subcenters"
            element={<CenterOwnerSubCenters />}
          />
          <Route
            path="/center-owner/teachers"
            element={<CenterOwnerTeachers />}
          />
          <Route
            path="/center-owner/teachers/:slug"
            element={<CenterOwnerTeacherDetail />}
          />

          {/* Center Scanner */}
          <Route
            path="/center-scanner/scanner"
            element={<CenterScannerScanner />}
          />
          <Route
            path="/center-scanner/students"
            element={<CenterScannerStudents />}
          />

          {/* Center Receptionist */}
          <Route
            path="/center-receptionist/dashboard"
            element={
              <Dashboard dashboardHref="/center-receptionist/dashboard" />
            }
          />
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </Suspense>
  )
}

export default App
