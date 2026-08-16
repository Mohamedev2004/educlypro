import { LoaderCircle } from "lucide-react"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import type { Staff } from "@/api/types/staff.types"

/**
 * Delete confirmation for a staff member.
 *
 * Responsibility: confirm the destructive action before it fires.
 * Layer: Components (domain)
 */

interface DeleteStaffDialogProps {
  t: (key: string) => string
  staff: Staff | null
  onOpenChange: (open: boolean) => void
  onConfirm: () => void
  isDeleting: boolean
}

export function DeleteStaffDialog({
  t,
  staff,
  onOpenChange,
  onConfirm,
  isDeleting,
}: DeleteStaffDialogProps) {
  return (
    <AlertDialog open={!!staff} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t("staff.deleteDialog.title")}</AlertDialogTitle>
          <AlertDialogDescription>
            {staff
              ? t("staff.deleteDialog.description").replace("{name}", staff.username)
              : ""}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isDeleting}>{t("common.cancel")}</AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            onClick={(event) => {
              event.preventDefault()
              onConfirm()
            }}
            disabled={isDeleting}
          >
            {isDeleting && <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />}
            {t("common.delete")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
