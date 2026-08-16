/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-refresh/only-export-components */
/* eslint-disable react-hooks/exhaustive-deps */
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react"
import { ar } from "../../locales/ar"
import { en } from "../../locales/en"
import { fr } from "../../locales/fr"
import { DirectionProvider } from "../../components/ui/direction"

type Direction = "ltr" | "rtl"
export type Locale = "ar" | "en" | "fr"

type TranslationTree = any

type DirectionContextValue = {
  direction: Direction
  locale: Locale
  toggleDirection: () => void
  setDirection: (dir: Direction) => void
  toggleLocale: () => void
  setLocale: (locale: Locale) => void
  t: (key: string, fallback?: string) => any
}

const messages: Record<Locale, any> = {
  ar,
  en,
  fr,
}

const DirectionContext = createContext<DirectionContextValue>({
  direction: "rtl",
  locale: "ar",
  toggleDirection: () => {},
  setDirection: () => {},
  toggleLocale: () => {},
  setLocale: () => {},
  t: (key, fallback) => (fallback ?? key),
})

function getInitialLocale(): Locale {
  if (typeof window === "undefined") {
    return "ar"
  }

  const savedLocale = getCookie("locale")
  if (savedLocale === "ar" || savedLocale === "en" || savedLocale === "fr") {
    return savedLocale
  }

  const legacyDirection = getCookie("direction")
  if (legacyDirection === "rtl") return "ar"
  if (legacyDirection === "ltr") return "en"

  return "ar"
}

function resolveMessage(tree: TranslationTree, key: string): any {
  const value = key.split(".").reduce((current, segment) => {
    if (!current || typeof current === "string" || Array.isArray(current)) return undefined
    return current[segment]
  }, tree)

  return typeof value === "string" || Array.isArray(value) ? value : undefined
}

function setCookie(name: string, value: string, days = 365) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${value}; expires=${expires}; path=/`
}

function getCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp("(^| )" + name + "=([^;]+)"))
  return match ? match[2] : null
}

export function DirectionProviderWrapper({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(getInitialLocale)
  const direction: Direction = locale === "ar" ? "rtl" : "ltr"

  const setLocale = useCallback((nextLocale: Locale) => {
    setLocaleState(nextLocale)
    setCookie("locale", nextLocale)
    setCookie("direction", nextLocale === "ar" ? "rtl" : "ltr")
  }, [])

  const setDirection = useCallback((dir: Direction) => {
    setLocale(dir === "rtl" ? "ar" : "en")
  }, [setLocale])

  const toggleDirection = () => setDirection(direction === "rtl" ? "ltr" : "rtl")
  const toggleLocale = () =>
    setLocale(locale === "ar" ? "en" : locale === "en" ? "fr" : "ar")

  useEffect(() => {
    document.documentElement.setAttribute("dir", direction)
    document.documentElement.setAttribute("lang", locale)
    const title = resolveMessage(messages[locale], "shell.appName")
    document.title = (typeof title === "string" ? title : "Lumina")
  }, [direction, locale])

  const value = useMemo<DirectionContextValue>(
    () => ({
      direction,
      locale,
      toggleDirection,
      setDirection,
      toggleLocale,
      setLocale,
      t: (key: string, fallback?: string) => {
        return resolveMessage(messages[locale], key) ?? (fallback || key)
      },
    }),
    [direction, locale, setDirection, toggleDirection, toggleLocale, setLocale],
  )

  return (
    <DirectionContext.Provider value={value}>
      <DirectionProvider dir={direction}>{children}</DirectionProvider>
    </DirectionContext.Provider>
  )
}

export const useDirection = () => useContext(DirectionContext)
export const useI18n = () => useContext(DirectionContext)
