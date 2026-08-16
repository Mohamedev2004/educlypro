import AppLayout from "../../layouts/app-layout"
import SettingsLayout from "../../layouts/setting-layout"
import { Input } from "../../components/ui/input"
import { Label } from "../../components/ui/label"
import { Button } from "../../components/ui/button"
import { useDirection } from "../../context/direction/direction-provider"
import { useUpdateProfile } from "../../hooks/auth/use-update-profile"

export default function Profile() {
  const { t } = useDirection()
  const {
    username,
    setUsername,
    email,
    setEmail,
    isSubmitting,
    errors,
    handleUpdateProfile,
  } = useUpdateProfile(t)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await handleUpdateProfile({ username, email })
  }

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("shell.settings") },
        { label: t("nav.profile") },
      ]}
    >
      <SettingsLayout>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <h2 className="text-lg font-medium">
              {t("settings.profileTitle")}
            </h2>
            <p className="text-sm text-muted-foreground">
              {t("settings.profileDescription")}
            </p>
          </div>

          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="username">{t("settings.name")}</Label>
              <Input
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder={t("settings.yourName")}
                className={errors.username ? "border-destructive" : ""}
                required
              />
              {errors.username && (
                <span className="text-xs text-destructive">
                  {errors.username}
                </span>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="email">{t("auth.email")}</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder={t("auth.emailPlaceholder")}
                className={errors.email ? "border-destructive" : ""}
                required
              />
              {errors.email && (
                <span className="text-xs text-destructive">{errors.email}</span>
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
