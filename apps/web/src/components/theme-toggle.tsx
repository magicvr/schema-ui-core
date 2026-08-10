import { Moon, Sun } from "lucide-react";

import { Button } from "@/components/ui/button";
import { setTheme } from "@/theme/theme";

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
        // setTheme persists the choice and applies it via applyThemeToElement,
        // which also syncs CSS color-scheme so native controls match (D8).
        setTheme(isDark() ? "light" : "dark");
      }}
    >
      <Sun className="h-4 w-4 dark:hidden" />
      <Moon className="hidden h-4 w-4 dark:block" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
