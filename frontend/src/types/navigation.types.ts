import type { LucideIcon } from "lucide-react"

/**
 * Navigation related types.
 * 
 * Responsibility: Define interfaces for navigation items and groups.
 * Layer: Types
 */

export type RolePageItem = {
  title: string
  url: string
  icon: LucideIcon
}

export type RolePageGroup = {
  label: string
  items: RolePageItem[]
}
