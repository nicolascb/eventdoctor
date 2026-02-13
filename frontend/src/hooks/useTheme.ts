import { useCallback, useEffect, useState } from "react";

type Theme = "light" | "dark";

function getStoredTheme(): Theme {
    const stored = localStorage.getItem("eventdoctor-theme") as Theme | null;
    if (stored === "light" || stored === "dark") return stored;

    // Respect OS preference
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) return "dark";
    return "light";
}

export function useTheme() {
    const [theme, setThemeState] = useState<Theme>(getStoredTheme);

    const applyTheme = useCallback((t: Theme) => {
        document.documentElement.classList.toggle("dark", t === "dark");
    }, []);

    useEffect(() => {
        applyTheme(theme);
    }, [theme, applyTheme]);

    // Listen for OS preference changes
    useEffect(() => {
        const mq = window.matchMedia("(prefers-color-scheme: dark)");
        const handler = (e: MediaQueryListEvent) => {
            if (!localStorage.getItem("eventdoctor-theme")) {
                const next = e.matches ? "dark" : "light";
                setThemeState(next);
            }
        };
        mq.addEventListener("change", handler);
        return () => mq.removeEventListener("change", handler);
    }, []);

    const toggleTheme = useCallback(() => {
        setThemeState((prev) => {
            const next = prev === "light" ? "dark" : "light";
            localStorage.setItem("eventdoctor-theme", next);
            return next;
        });
    }, []);

    return { theme, toggleTheme };
}
