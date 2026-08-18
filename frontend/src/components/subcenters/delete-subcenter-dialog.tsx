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
import type { SubCenter } from "@/api/types/subcenters.types"

/**
 * Delete confirmation for a sub-center.
 *
 * Responsibility: confirm the destructive action before it fires. The
 * backend refuses the delete (and this surfaces as a toast) if staff are
 * still assigned to the sub-center.
 * Layer: Components (domain)
 */

interface DeleteSubCenterDialogProps {
  t: (key: string) => string
  subCenter: SubCenter | null
  onOpenChange: (open: boolean) => void
  onConfirm: () => void
  isDeleting: boolean
}

export function DeleteSubCenterDialog({
  t,
  subCenter,
  onOpenChange,
  onConfirm,
  isDeleting,
}: DeleteSubCenterDialogProps) {
  return (
    <AlertDialog open={!!subCenter} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            {t("subcenters.deleteDialog.title")}
          </AlertDialogTitle>
          <AlertDialogDescription>
            {subCenter
              ? t("subcenters.deleteDialog.description").replace(
                  "{name}",
                  subCenter.name
                )
              : ""}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isDeleting}>
            {t("common.cancel")}
          </AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            onClick={(event) => {
              event.preventDefault()
              onConfirm()
            }}
            disabled={isDeleting}
          >
            {isDeleting && (
              <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />
            )}
            {t("common.delete")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
