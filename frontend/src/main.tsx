import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter } from "react-router-dom"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"

import "./index.css"
import App from "./App.tsx"
import { ThemeProvider } from "./components/theme-provider.tsx"
import { TooltipProvider } from "./components/ui/tooltip.tsx"
import { DirectionProviderWrapper } from "./context/direction/direction-provider.tsx"
import { AuthProvider } from "./context/auth/auth-context.tsx"
import { CommandPaletteProvider } from "./context/command-palette/command-palette-context.tsx"
import { AppToaster } from "./components/app-toaster.tsx"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      refetchOnReconnect: false,
      retry: false,
    },
  },
})

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <TooltipProvider>
          <BrowserRouter>
            <DirectionProviderWrapper>
              <AuthProvider>
                <CommandPaletteProvider>
                  <App />
                  <AppToaster richColors/>
                </CommandPaletteProvider>
              </AuthProvider>
            </DirectionProviderWrapper>
          </BrowserRouter>
        </TooltipProvider>
      </ThemeProvider>
    </QueryClientProvider>
  </StrictMode>
)
