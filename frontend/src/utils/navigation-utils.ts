/* eslint-disable @typescript-eslint/no-explicit-any */
import type { Role, User } from "@/api/types/auth.types"
import type { RolePageGroup } from "@/types/navigation.types"
import {
  BellRing,
  University,
  ChartBar,
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

export function getDashboardPath(role: Role): string {
  return dashboardPathByRole[role]
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
            title: t("roles.staff"),
            url: "/center-owner/staff",
            icon: UserKey,
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
