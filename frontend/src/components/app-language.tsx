import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../components/ui/select"

export function AppLanguage({
  locale,
  setLocale,
}: {
  locale: string
  setLocale: (value: string) => void
}) {
  return (
    <Select value={locale} onValueChange={setLocale}>
      <SelectTrigger className="!h-9 sm:w-[160px] max-w-max">
        <SelectValue>
          {locale === "en" && "English 🇬🇧"}
          {locale === "fr" && "Français 🇫🇷"}
          {locale === "ar" && "العربية 🇲🇦"}
        </SelectValue>
      </SelectTrigger>

      <SelectContent>
        <SelectItem value="en">English 🇬🇧</SelectItem>
        <SelectItem value="fr">Français 🇫🇷</SelectItem>
        <SelectItem value="ar">العربية 🇲🇦</SelectItem>
      </SelectContent>
    </Select>
  )
}
