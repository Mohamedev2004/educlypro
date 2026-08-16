import { Toaster as Sonner, type ToasterProps } from "sonner"
import { useDirection } from "@/context/direction/direction-provider"

export function AppToaster(props: ToasterProps) {
  const { direction } = useDirection()

  const position = direction === "rtl" ? "top-center" : "top-center"

  return (
    <Sonner
      position={position}
      duration={3000}
      visibleToasts={6}
      className="toaster"
      icons={{
        success: null,
        info: null,
        warning: null,
        error: null,
        loading: null,
      }}
      style={
        {
          /* General & Default Toasts */
          "--normal-bg": "var(--background)",
          "--normal-text": "var(--foreground)",
          "--normal-border": "var(--border)",
          "--border-radius": "var(--radius)",

          /* Success - Monochrome High Contrast */
          "--success-bg": "var(--primary)",        
          "--success-text": "var(--primary-foreground)",
          "--success-border": "var(--primary)",  
          
          /* Error - Using your theme's destructive variable */
          "--error-bg": "color-mix(in oklch, var(--destructive) 10%, transparent)",
          "--error-text": "var(--destructive)",
          "--error-border": "color-mix(in oklch, var(--destructive) 10%, transparent)",

          /* Info/Warning - keeping it clean */
          "--info-bg": "var(--background)",
          "--info-text": "var(--foreground)",
          "--warning-text": "var(--foreground)",
        } as React.CSSProperties
      }
      {...props}
    />
  )
}