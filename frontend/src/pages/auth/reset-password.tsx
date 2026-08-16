import { useState } from "react"
import { Link } from "react-router-dom"

import AuthLayout from "../../layouts/auth-layout"
import { Button } from "../../components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../../components/ui/card"
import { PasswordInput } from "../../components/ui/password-input"
import { Label } from "../../components/ui/label"
import { useDirection } from "../../context/direction/direction-provider"
import { useResetPassword } from "../../hooks/auth/use-reset-password"

export default function ResetPassword() {
  const { t } = useDirection()
  const { isSubmitting, errors, handleResetPassword } = useResetPassword(t)

  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await handleResetPassword({ password, confirmPassword })
  }

  return (
    <AuthLayout>
      <div className="flex justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle>{t("auth.resetPasswordTitle")}</CardTitle>
            <CardDescription>
              {t("auth.resetPasswordDescription")}
            </CardDescription>
            <CardAction>
              <Link to="/login">
                <Button variant="link" className="text-foreground">
                  {t("auth.logIn")}
                </Button>
              </Link>
            </CardAction>
          </CardHeader>

          <CardContent>
            <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
              <div className="grid gap-2">
                <Label htmlFor="password">{t("auth.newPassword")}</Label>
                <PasswordInput
                  id="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  className={errors.password ? "border-destructive" : ""}
                  required
                />
                {errors.password && (
                  <span className="text-xs text-destructive">
                    {errors.password}
                  </span>
                )}
              </div>

              <div className="grid gap-2">
                <Label htmlFor="confirmPassword">
                  {t("auth.confirmNewPassword")}
                </Label>
                <PasswordInput
                  id="confirmPassword"
                  value={confirmPassword}
                  onChange={(event) => setConfirmPassword(event.target.value)}
                  className={errors.confirmPassword ? "border-destructive" : ""}
                  required
                />
                {errors.confirmPassword && (
                  <span className="text-xs text-destructive">
                    {errors.confirmPassword}
                  </span>
                )}
              </div>

              {errors.general && (
                <div className="text-sm font-medium text-destructive">
                  {errors.general}
                </div>
              )}

              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting
                  ? t("auth.resettingPassword")
                  : t("auth.resetPassword")}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  )
}
