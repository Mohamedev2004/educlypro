import type { FormEventHandler } from "react"
import { LoaderCircle } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { PasswordInput } from "@/components/ui/password-input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useStaffForm } from "@/hooks/staff/use-staff-form"
import type { Staff } from "@/api/types/staff.types"

/**
 * Create/edit dialog for a staff member.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useStaffForm`. `centerId` is a fixed prop (the signed-in owner's own
 * center) — there is no center picker; it's always included in the
 * submitted payload, never user-editable.
 * Layer: Components (domain)
 */

interface StaffFormDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  centerId: number
  staff: Staff | null
  onSuccess: () => void
}

export function StaffFormDialog({
  t,
  open,
  onOpenChange,
  centerId,
  staff,
  onSuccess,
}: StaffFormDialogProps) {
  const {
    isEdit,
    username,
    setUsername,
    email,
    setEmail,
    role,
    setRole,
    password,
    setPassword,
    errors,
    isSubmitting,
    handleSubmit,
  } = useStaffForm(t, centerId, staff, () => {
    onOpenChange(false)
    onSuccess()
  })

  const submit: FormEventHandler = async (event) => {
    event.preventDefault()
    await handleSubmit()
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {isEdit ? t("staff.dialog.editTitle") : t("staff.dialog.addTitle")}
          </DialogTitle>
          <DialogDescription>
            {isEdit ? t("staff.dialog.editDescription") : t("staff.dialog.addDescription")}
          </DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="staff-username">{t("staff.fields.username")}</Label>
            <Input
              id="staff-username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={errors.username ? "border-destructive" : ""}
              autoFocus
              required
            />
            {errors.username && (
              <span className="text-xs text-destructive">{errors.username}</span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="staff-email">{t("staff.fields.email")}</Label>
            <Input
              id="staff-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={errors.email ? "border-destructive" : ""}
              required
            />
            {errors.email && (
              <span className="text-xs text-destructive">{errors.email}</span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="staff-role">{t("staff.fields.role")}</Label>
            <Select value={role} onValueChange={(value) => setRole(value as typeof role)}>
              <SelectTrigger id="staff-role" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="center_scanner">{t("roles.center_scanner")}</SelectItem>
                <SelectItem value="center_receptionist">
                  {t("roles.center_receptionist")}
                </SelectItem>
              </SelectContent>
            </Select>
            {errors.role && <span className="text-xs text-destructive">{errors.role}</span>}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="staff-password">
              {t("staff.fields.password")}
              {isEdit && (
                <span className="text-muted-foreground"> {t("staff.fields.passwordEditHint")}</span>
              )}
            </Label>
            <PasswordInput
              id="staff-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={errors.password ? "border-destructive" : ""}
              autoComplete="new-password"
              required={!isEdit}
            />
            {errors.password && (
              <span className="text-xs text-destructive">{errors.password}</span>
            )}
          </div>

          {errors.general && (
            <div className="text-sm font-medium text-destructive">{errors.general}</div>
          )}

          <DialogFooter>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />}
              {isEdit ? t("staff.dialog.submitEdit") : t("staff.dialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
