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
import { useAddStaffForm } from "@/hooks/centers/use-add-staff-form"
import { useCenterSubCenters } from "@/hooks/centers/use-center-subcenters"

/**
 * Add-staff dialog for the center detail page.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useAddStaffForm`.
 * Layer: Components (domain)
 */

interface AddStaffDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  centerSlug: string
  onSuccess: () => void
}

export function AddStaffDialog({
  t,
  open,
  onOpenChange,
  centerSlug,
  onSuccess,
}: AddStaffDialogProps) {
  const {
    username,
    setUsername,
    email,
    setEmail,
    role,
    setRole,
    password,
    setPassword,
    subCenterId,
    setSubCenterId,
    errors,
    isSubmitting,
    handleSubmit,
  } = useAddStaffForm(t, centerSlug, () => {
    onOpenChange(false)
    onSuccess()
  })

  const { data: subCentersData } = useCenterSubCenters(centerSlug)
  const subCenters = subCentersData?.items ?? []

  const submit: FormEventHandler = async (event) => {
    event.preventDefault()
    await handleSubmit()
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("centers.addStaffDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("centers.addStaffDialog.description")}
          </DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="center-staff-username">
              {t("centers.fields.username")}
            </Label>
            <Input
              id="center-staff-username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={errors.username ? "border-destructive" : ""}
              autoFocus
              required
            />
            {errors.username && (
              <span className="text-xs text-destructive">
                {errors.username}
              </span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="center-staff-email">
              {t("centers.fields.email")}
            </Label>
            <Input
              id="center-staff-email"
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
            <Label htmlFor="center-staff-role">
              {t("centers.fields.role")}
            </Label>
            <Select
              value={role}
              onValueChange={(value) => setRole(value as typeof role)}
            >
              <SelectTrigger id="center-staff-role" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="center_scanner">
                  {t("roles.center_scanner")}
                </SelectItem>
                <SelectItem value="center_receptionist">
                  {t("roles.center_receptionist")}
                </SelectItem>
              </SelectContent>
            </Select>
            {errors.role && (
              <span className="text-xs text-destructive">{errors.role}</span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="center-staff-subcenter">
              {t("centers.fields.subCenter")}
            </Label>
            <Select
              value={subCenterId ? String(subCenterId) : ""}
              onValueChange={(value) => setSubCenterId(Number(value))}
            >
              <SelectTrigger
                id="center-staff-subcenter"
                className={
                  errors.subCenterId ? "w-full border-destructive" : "w-full"
                }
              >
                <SelectValue
                  placeholder={t("centers.fields.subCenterPlaceholder")}
                />
              </SelectTrigger>
              <SelectContent>
                {subCenters.map((subCenter) => (
                  <SelectItem key={subCenter.id} value={String(subCenter.id)}>
                    {subCenter.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.subCenterId && (
              <span className="text-xs text-destructive">
                {errors.subCenterId}
              </span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="center-staff-password">
              {t("centers.fields.password")}
            </Label>
            <PasswordInput
              id="center-staff-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={errors.password ? "border-destructive" : ""}
              autoComplete="new-password"
              required
            />
            {errors.password && (
              <span className="text-xs text-destructive">
                {errors.password}
              </span>
            )}
          </div>

          {errors.general && (
            <div className="text-sm font-medium text-destructive">
              {errors.general}
            </div>
          )}

          <DialogFooter>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && (
                <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />
              )}
              {t("centers.addStaffDialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
