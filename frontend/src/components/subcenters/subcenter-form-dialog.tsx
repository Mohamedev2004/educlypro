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
import { useSubCenterForm } from "@/hooks/subcenters/use-subcenter-form"
import type { SubCenter } from "@/api/types/subcenters.types"

/**
 * Create/edit dialog for a sub-center.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useSubCenterForm`.
 * Layer: Components (domain)
 */

interface SubCenterFormDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  subCenter: SubCenter | null
  onSuccess: () => void
}

export function SubCenterFormDialog({
  t,
  open,
  onOpenChange,
  subCenter,
  onSuccess,
}: SubCenterFormDialogProps) {
  const { isEdit, name, setName, errors, isSubmitting, handleSubmit } =
    useSubCenterForm(t, open, subCenter, () => {
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
              ? t("subcenters.dialog.editTitle")
              : t("subcenters.dialog.addTitle")}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t("subcenters.dialog.editDescription")
              : t("subcenters.dialog.addDescription")}
          </DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="subcenter-name">
              {t("subcenters.fields.name")}
            </Label>
            <Input
              id="subcenter-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className={errors.name ? "border-destructive" : ""}
              autoFocus
              required
            />
            {errors.name && (
              <span className="text-xs text-destructive">{errors.name}</span>
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
                ? t("subcenters.dialog.submitEdit")
                : t("subcenters.dialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
