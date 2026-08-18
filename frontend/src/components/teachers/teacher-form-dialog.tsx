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
import { Label } from "@/components/ui/label"
import { useTeacherForm } from "@/hooks/teachers/use-teacher-form"
import type { Teacher } from "@/api/types/teachers.types"

/**
 * Create/edit dialog for a teacher's core identity fields (full name,
 * email, phone). Subject/class assignment happens from the teacher detail
 * page instead.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useTeacherForm`.
 * Layer: Components (domain)
 */

interface TeacherFormDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  teacher: Teacher | null
  onSuccess: () => void
}

export function TeacherFormDialog({
  t,
  open,
  onOpenChange,
  teacher,
  onSuccess,
}: TeacherFormDialogProps) {
  const {
    isEdit,
    fullName,
    setFullName,
    email,
    setEmail,
    phone,
    setPhone,
    errors,
    isSubmitting,
    handleSubmit,
  } = useTeacherForm(t, open, teacher, () => {
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
            {isEdit
              ? t("teachers.dialog.editTitle")
              : t("teachers.dialog.addTitle")}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t("teachers.dialog.editDescription")
              : t("teachers.dialog.addDescription")}
          </DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="teacher-full-name">
              {t("teachers.fields.fullName")}
            </Label>
            <Input
              id="teacher-full-name"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              className={errors.fullName ? "border-destructive" : ""}
              autoFocus
              required
            />
            {errors.fullName && (
              <span className="text-xs text-destructive">
                {errors.fullName}
              </span>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="teacher-email">{t("teachers.fields.email")}</Label>
            <Input
              id="teacher-email"
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
            <Label htmlFor="teacher-phone">{t("teachers.fields.phone")}</Label>
            <Input
              id="teacher-phone"
              type="tel"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              className={errors.phone ? "border-destructive" : ""}
              required
            />
            {errors.phone && (
              <span className="text-xs text-destructive">{errors.phone}</span>
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
              {isEdit
                ? t("teachers.dialog.submitEdit")
                : t("teachers.dialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
