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
import type { Teacher } from "@/api/types/teachers.types"

/**
 * Delete confirmation for a teacher.
 *
 * Responsibility: confirm the destructive action before it fires.
 * Layer: Components (domain)
 */

interface DeleteTeacherDialogProps {
  t: (key: string) => string
  teacher: Teacher | null
  onOpenChange: (open: boolean) => void
  onConfirm: () => void
  isDeleting: boolean
}

export function DeleteTeacherDialog({
  t,
  teacher,
  onOpenChange,
  onConfirm,
  isDeleting,
}: DeleteTeacherDialogProps) {
  return (
    <AlertDialog open={!!teacher} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            {t("teachers.deleteDialog.title")}
          </AlertDialogTitle>
          <AlertDialogDescription>
            {teacher
              ? t("teachers.deleteDialog.description").replace(
                  "{name}",
                  teacher.full_name
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
