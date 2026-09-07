import * as React from "react";
import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ theme, ...props }: ToasterProps) => {
  const [currentTheme, setCurrentTheme] = React.useState<ToasterProps["theme"]>(theme || "system");

  React.useEffect(() => {
    if (theme) {
      setCurrentTheme(theme);
      return;
    }
    const isDark = typeof document !== "undefined" && document.documentElement.classList.contains("dark");
    setCurrentTheme(isDark ? "dark" : "light");
  }, [theme]);

  return (
    <Sonner
      theme={currentTheme}
      className="toaster group"
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
        } as React.CSSProperties
      }
      {...props}
    />
  );
};

export { Toaster };
