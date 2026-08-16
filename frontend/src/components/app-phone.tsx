import { Phone } from "lucide-react";
import { Button } from "./ui/button";

export function AppPhone() {
  return (
    <a href="tel:+212695645950">
      <Button
        variant="default"
        className="size-9"
      >
        <Phone className="size-4" />
      </Button>
    </a>
  )
}