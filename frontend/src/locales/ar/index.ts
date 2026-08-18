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
import { centers } from "./centers"
import { table } from "./table"
import { onboarding } from "./onboarding"
import { subcenters } from "./subcenters"
import { teachers } from "./teachers"

export const ar = {
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
  centers,
  table,
  onboarding,
  subcenters,
  teachers,
} as const
