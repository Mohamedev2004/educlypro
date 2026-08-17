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
import { useCenterForm } from "@/hooks/centers/use-center-form"

/**
 * Create dialog for a center.
 *
 * Responsibility: render the form UI, delegate all state/submission to
 * `useCenterForm`.
 * Layer: Components (domain)
 */

interface CenterFormDialogProps {
  t: (key: string) => string
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function CenterFormDialog({ t, open, onOpenChange, onSuccess }: CenterFormDialogProps) {
  const { name, setName, errors, isSubmitting, handleSubmit } = useCenterForm(t, () => {
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
          <DialogTitle>{t("centers.dialog.addTitle")}</DialogTitle>
          <DialogDescription>{t("centers.dialog.addDescription")}</DialogDescription>
        </DialogHeader>

        <form className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="center-name">{t("centers.fields.name")}</Label>
            <Input
              id="center-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className={errors.name ? "border-destructive" : ""}
              autoFocus
              required
            />
            {errors.name && <span className="text-xs text-destructive">{errors.name}</span>}
          </div>

          {errors.general && (
            <div className="text-sm font-medium text-destructive">{errors.general}</div>
          )}

          <DialogFooter>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />}
              {t("centers.dialog.submitAdd")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
