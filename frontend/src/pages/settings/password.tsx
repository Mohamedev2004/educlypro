import { useState } from "react"
import AppLayout from "../../layouts/app-layout"
import SettingsLayout from "../../layouts/setting-layout"
import { PasswordInput } from "../../components/ui/password-input"
import { Label } from "../../components/ui/label"
import { Button } from "../../components/ui/button"
import { useDirection } from "../../context/direction/direction-provider"
import { useUpdatePassword } from "../../hooks/auth/use-update-password"

export default function Password() {
  const { t } = useDirection()
  const { isSubmitting, errors, handleUpdatePassword } = useUpdatePassword(t)
  const [currentPassword, setCurrentPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const success = await handleUpdatePassword({
      currentPassword,
      newPassword,
      confirmPassword,
    })

    if (success) {
      setCurrentPassword("")
      setNewPassword("")
      setConfirmPassword("")
    }
  }

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("shell.settings") },
        { label: t("nav.password") },
      ]}
    >
      <SettingsLayout>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <h2 className="text-lg font-medium">
              {t("settings.passwordTitle")}
            </h2>
            <p className="text-sm text-muted-foreground">
              {t("settings.passwordDescription")}
            </p>
          </div>

          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="currentPassword">
                {t("settings.currentPassword")}
              </Label>
              <PasswordInput
                id="currentPassword"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                placeholder={t("settings.enterCurrentPassword")}
                className={errors.currentPassword ? "border-destructive" : ""}
                required
              />
              {errors.currentPassword && (
                <span className="text-xs text-destructive">
                  {errors.currentPassword}
                </span>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="newPassword">{t("settings.newPassword")}</Label>
              <PasswordInput
                id="newPassword"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder={t("settings.enterNewPassword")}
                className={errors.newPassword ? "border-destructive" : ""}
                required
              />
              {errors.newPassword && (
                <span className="text-xs text-destructive">
                  {errors.newPassword}
                </span>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="confirmPassword">
                {t("settings.confirmPassword")}
              </Label>
              <PasswordInput
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder={t("settings.confirmNewPassword")}
                className={errors.confirmPassword ? "border-destructive" : ""}
                required
              />
              {errors.confirmPassword && (
                <span className="text-xs text-destructive">
                  {errors.confirmPassword}
                </span>
              )}
            </div>
          </div>

          {errors.general && (
            <div className="text-sm font-medium text-destructive">
              {errors.general}
            </div>
          )}

          <div className="flex justify-start">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? t("common.saving") : t("settings.saveChanges")}
            </Button>
          </div>
        </form>
      </SettingsLayout>
    </AppLayout>
  )
}
