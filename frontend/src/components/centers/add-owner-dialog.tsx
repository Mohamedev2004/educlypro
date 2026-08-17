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
import { useAddOwnerForm } from "@/hooks/centers/use-add-owner-form"

/**
 * Add-owner dialog for a center that doesn't have one yet.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useAddOwnerForm`.
 * Layer: Components (domain)
 */

interface AddOwnerDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  centerSlug: string
  onSuccess: () => void
}

export function AddOwnerDialog({
  t,
  open,
  onOpenChange,
  centerSlug,
  onSuccess,
}: AddOwnerDialogProps) {
  const {
    username,
    setUsername,
    email,
    setEmail,
    password,
    setPassword,
    errors,
    isSubmitting,
    handleSubmit,
  } = useAddOwnerForm(t, centerSlug, () => {
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
          <DialogTitle>{t("centers.ownerDialog.title")}</DialogTitle>
          <DialogDescription>{t("centers.ownerDialog.description")}</DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="owner-username">{t("centers.fields.username")}</Label>
            <Input
              id="owner-username"
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
            <Label htmlFor="owner-email">{t("centers.fields.email")}</Label>
            <Input
              id="owner-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={errors.email ? "border-destructive" : ""}
              required
            />
            {errors.email && <span className="text-xs text-destructive">{errors.email}</span>}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="owner-password">{t("centers.fields.password")}</Label>
            <PasswordInput
              id="owner-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={errors.password ? "border-destructive" : ""}
              autoComplete="new-password"
              required
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
              {t("centers.ownerDialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
