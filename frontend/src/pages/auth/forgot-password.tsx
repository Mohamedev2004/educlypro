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
import { Input } from "../../components/ui/input"
import { Label } from "../../components/ui/label"
import { useDirection } from "../../context/direction/direction-provider"
import { useForgotPassword } from "../../hooks/auth/use-forgot-password"

export default function ForgotPassword() {
  const { t } = useDirection()
  const { isSubmitting, errors, handleForgotPassword } = useForgotPassword(t)
  const [email, setEmail] = useState("")

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await handleForgotPassword({ email })
  }

  return (
    <AuthLayout>
      <div className="flex justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle>{t("auth.forgotPasswordTitle")}</CardTitle>
            <CardDescription>
              {t("auth.forgotPasswordDescription")}
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
                <Label htmlFor="email">{t("auth.email")}</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder={t("auth.emailPlaceholder")}
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  className={errors.email ? "border-destructive" : ""}
                  required
                />
                {errors.email && (
                  <span className="text-xs text-destructive">
                    {errors.email}
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
                  ? t("auth.sendingResetLink")
                  : t("auth.sendResetLink")}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  )
}
