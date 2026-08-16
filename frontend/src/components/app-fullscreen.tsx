import { useEffect, useState } from "react"
import { Button } from "./ui/button"
import { Maximize, Minimize } from "lucide-react"

export function AppFullscreen() {
  const [isFullscreen, setIsFullscreen] = useState(
    () => !!document.fullscreenElement  // initialize from DOM, not false
  )

  useEffect(() => {
    const handleChange = () => {
      const fullscreen = !!document.fullscreenElement
      setIsFullscreen(fullscreen)
      document.documentElement.classList.toggle("hide-scrollbar", fullscreen)
    }

    document.addEventListener("fullscreenchange", handleChange)
    return () => document.removeEventListener("fullscreenchange", handleChange)
  }, [])

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen()
    } else {
      document.exitFullscreen()
    }
  }

  return (
    <div className="hidden sm:block">
      <Button
        variant="default"
        size="icon"
        onClick={toggleFullscreen}
        className="size-9"
      >
        {isFullscreen ? (
          <Minimize className="size-4" />
        ) : (
          <Maximize className="size-4" />
        )}
        <span className="sr-only">Toggle fullscreen</span>
      </Button>
    </div>
  )
}
