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
import { PasswordInput } from "../../components/ui/password-input"
import { Label } from "../../components/ui/label"
import { useDirection } from "../../context/direction/direction-provider"
import { useLogin } from "../../hooks/auth/use-login"

export default function Login() {
  const { t } = useDirection()
  const { isSubmitting, errors, handleLogin } = useLogin(t)
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    await handleLogin({ email, password })
  }

  return (
    <AuthLayout>
      <div className="flex justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle>{t("auth.loginTitle")}</CardTitle>
            <CardDescription>{t("auth.loginDescription")}</CardDescription>
            <CardAction>
              {/* <Link to="/register">
                <Button variant="link" className="text-foreground">
                  {t("auth.signUp")}
                </Button>
              </Link> */}
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

              <div className="grid gap-2">
                <div className="flex items-center">
                  <Label htmlFor="password">{t("auth.password")}</Label>
                  <Link
                    to="/forgot-password"
                    className="ms-auto text-sm underline hover:underline"
                  >
                    {t("auth.forgotPassword")}
                  </Link>
                </div>
                <PasswordInput
                  id="password"
                  placeholder={t("auth.password")}
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

              {errors.general && (
                <div className="text-sm font-medium text-destructive">
                  {errors.general}
                </div>
              )}

              <Button className="w-full" type="submit" disabled={isSubmitting}>
                {isSubmitting ? t("auth.loggingIn") : t("auth.login")}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  )
}
