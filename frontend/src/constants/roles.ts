import { type Role } from "@/api/types/auth.types"

/**
 * Role related constants.
 *
 * Responsibility: Static values for role-based labels and configurations.
 * Layer: Constants
 */

export const roleLabels: Record<Role, string> = {
  super_admin: "roles.super_admin",
  center_owner: "roles.center_owner",
  center_scanner: "roles.center_scanner",
  center_receptionist: "roles.center_receptionist",
}
