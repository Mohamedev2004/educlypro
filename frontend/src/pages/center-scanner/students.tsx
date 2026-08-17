import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { PlaceholderPattern } from "@/components/ui/placeholder-pattern"

export default function CenterScannerStudents() {
  const { t } = useDirection()

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.checkIn"), href: "/center-scanner/scanner" },
        { label: t("roles.students") },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <div className="grid auto-rows-min gap-4 md:grid-cols-3">
          <PlaceholderPattern className="aspect-video rounded-xl bg-sidebar" />
          <PlaceholderPattern className="aspect-video rounded-xl bg-sidebar" />
          <PlaceholderPattern className="aspect-video rounded-xl bg-sidebar" />
        </div>
        <PlaceholderPattern className="min-h-[100vh] flex-1 rounded-xl bg-sidebar" />
      </div>
    </AppLayout>
  )
}
