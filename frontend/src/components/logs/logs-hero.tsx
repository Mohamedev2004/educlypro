interface LogsHeroProps {
  t: (key: string, fallback?: string) => string
}

export function LogsHero({ t }: LogsHeroProps) {
  return (
    <div className="space-y-2">
      <h1 className="text-2xl font-semibold text-foreground">
        {t("logs.title")}
      </h1>

      <p className="text-sm text-muted-foreground max-w-3xl text-justify">
        {t("logs.description")}
      </p>
    </div>
  )
}