import { Moon, Sun } from "lucide-react";

import { Button } from "@/components/ui/button";

function isDark() {
  return document.documentElement.classList.contains("dark");
}

export function ThemeToggle() {
  return (
    <Button
      type="button"
      variant="outline"
      size="sm"
      aria-label="Toggle color theme"
      title="Toggle color theme"
      onClick={() => {
        const next = !isDark();
        document.documentElement.classList.toggle("dark", next);
        localStorage.setItem("theme", next ? "dark" : "light");
      }}
    >
      <Sun className="h-4 w-4 dark:hidden" />
      <Moon className="hidden h-4 w-4 dark:block" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
