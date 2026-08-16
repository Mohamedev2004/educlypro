import { shell } from "./shell"
import { nav } from "./nav"
import { common } from "./common"
import { authLayout } from "./auth-layout"
import { auth } from "./auth"
import { dashboard } from "./dashboard"
import { settings } from "./settings"
import { notFound } from "./not-found"
import { validation } from "./validation"
import { api } from "./api"
import { notifications } from "./notifications"
import { logs } from "./logs"
import { roles } from "./roles"
import { staff } from "./staff"
import { table } from "./table"

export const en = {
  shell,
  nav,
  common,
  authLayout,
  auth,
  dashboard,
  settings,
  notFound,
  validation,
  api,
  notifications,
  logs,
  roles,
  staff,
  table,
} as const
