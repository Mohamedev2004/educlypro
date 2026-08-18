/* eslint-disable @typescript-eslint/no-explicit-any */
import type { Role, User } from "@/api/types/auth.types"
import type { RolePageGroup } from "@/types/navigation.types"
import {
  BellRing,
  University,
  Building2,
  ChartBar,
  GraduationCap,
  History,
  QrCode,
  UserKey,
  Users,
} from "lucide-react"

/**
 * Navigation utility functions.
 *
 * Responsibility: Logic for generating role-based navigation groups and
 * resolving the single dashboard URL for a role.
 * Layer: Utils
 */

const dashboardPathByRole: Record<Role, string> = {
  super_admin: "/super-admin/dashboard",
  center_owner: "/center-owner/dashboard",
  center_scanner: "/center-scanner/scanner",
  center_receptionist: "/center-receptionist/dashboard",
}

// The one place a center_owner without grades is always allowed to be.
export const ONBOARDING_PATH = "/center-owner/onboarding"

export function getDashboardPath(role: Role): string {
  return dashboardPathByRole[role]
}

// A center_owner can't use the dashboard (or anything else past login)
// until their center has at least one grade — this is the "at least one
// grade" gate the whole app enforces.
export function needsOnboarding(user: User | null): boolean {
  return !!user && user.role === "center_owner" && !user.has_grades
}

// Where a user should land right after authenticating (login, or an
// already-authenticated visit to a public-only route): onboarding if the
// gate above applies, otherwise their role's dashboard.
export function getPostLoginPath(user: User | null): string {
  if (!user) return "/login"
  if (needsOnboarding(user)) return ONBOARDING_PATH
  return getDashboardPath(user.role)
}

export function getRolePages(
  user: User,
  t: (key: string) => any
): RolePageGroup[] {
  if (user.role === "super_admin") {
    return [
      {
        label: t("roles.overview"),
        items: [
          {
            title: t("roles.dashboard"),
            url: "/super-admin/dashboard",
            icon: ChartBar,
          },
          {
            title: t("roles.notifications"),
            url: "/super-admin/notifications",
            icon: BellRing,
          },
          {
            title: t("roles.logs"),
            url: "/super-admin/logs",
            icon: History,
          },
          {
            title: t("roles.centers"),
            url: "/super-admin/centers",
            icon: University,
          },
        ],
      },
    ]
  }

  if (user.role === "center_owner") {
    return [
      {
        label: t("roles.overview"),
        items: [
          {
            title: t("roles.dashboard"),
            url: getDashboardPath(user.role),
            icon: ChartBar,
          },
          {
            title: t("roles.notifications"),
            url: "/center-owner/notifications",
            icon: BellRing,
          },
          {
            title: t("roles.staff"),
            url: "/center-owner/staff",
            icon: UserKey,
          },
          {
            title: t("roles.subCenters"),
            url: "/center-owner/subcenters",
            icon: Building2,
          },
          {
            title: t("roles.teachers"),
            url: "/center-owner/teachers",
            icon: GraduationCap,
          },
        ],
      },
    ]
  }

  if (user.role === "center_scanner") {
    return [
      {
        label: t("roles.checkIn"),
        items: [
          {
            title: t("roles.scanner"),
            url: "/center-scanner/scanner",
            icon: QrCode,
          },
          {
            title: t("roles.students"),
            url: "/center-scanner/students",
            icon: Users,
          },
        ],
      },
    ]
  }

  return [
    {
      label: t("roles.overview"),
      items: [
        {
          title: t("roles.dashboard"),
          url: getDashboardPath(user.role),
          icon: ChartBar,
        },
      ],
    },
  ]
}
