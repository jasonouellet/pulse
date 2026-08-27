import React, { useState, useEffect } from "react";
import {
  Sun,
  Moon,
  Menu,
  X,
  Shield,
  Calendar,
  Users,
  Trophy,
  Settings,
} from "lucide-react";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

// ============================================================================
// 1. REUSABLE UI COMPONENTS (WCAG 2.1 AA Compliant)
// ============================================================================

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "outline";
  size?: "sm" | "md" | "lg";
}

const Button: React.FC<ButtonProps> = ({
  children,
  className,
  variant = "primary",
  size: _size = "md", // TODO: brancher les variantes de taille
  ...props
}) => {
  const baseStyles =
    "inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 min-h-[44px] min-w-[44px] px-4 py-2 disabled:opacity-50 disabled:pointer-events-none";

  const variants = {
    primary: "bg-brand-500 text-white hover:bg-brand-900",
    secondary:
      "bg-slate-200 text-slate-900 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-100",
    outline:
      "border border-slate-300 bg-transparent hover:bg-slate-100 dark:border-slate-600 dark:hover:bg-slate-800",
  };

  return (
    <button
      className={twMerge(clsx(baseStyles, variants[variant], className))}
      {...props}
    >
      {children}
    </button>
  );
};

// ============================================================================
// 2. LAYOUT COMPONENT (AppLayout)
// ============================================================================

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const [isDarkMode, setIsDarkMode] = useState<boolean>(() => {
    return (
      localStorage.getItem("theme") === "dark" ||
      (!("theme" in localStorage) &&
        window.matchMedia("(prefers-color-scheme: dark)").matches)
    );
  });

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState<boolean>(false);

  useEffect(() => {
    if (isDarkMode) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  }, [isDarkMode]);

  const toggleTheme = () => {
    setIsDarkMode((prev) => !prev);
  };

  const navItems = [
    { label: "Dashboard", href: "#", icon: Shield },
    { label: "Pools & Rosters", href: "#", icon: Users },
    { label: "Tournaments", href: "#", icon: Trophy },
    { label: "Schedule", href: "#", icon: Calendar },
    { label: "Settings", href: "#", icon: Settings },
  ];

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 transition-colors duration-200 dark:bg-slate-900 dark:text-slate-100">
      {/* Navigation Bar */}
      <header className="sticky top-0 z-40 w-full border-b border-slate-200 bg-white/80 backdrop-blur-md dark:border-slate-800 dark:bg-slate-900/80">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          {/* Logo Brand */}
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-500 font-bold text-white shadow-md shadow-brand-500/20">
              P
            </div>
            <span className="text-xl font-bold tracking-tight text-slate-900 dark:text-white">
              PULSE<span className="text-brand-500">.os</span>
            </span>
          </div>

          {/* Desktop Links */}
          <nav
            className="hidden md:flex md:items-center md:gap-1"
            aria-label="Main Navigation"
          >
            {navItems.map((item) => {
              const Icon = item.icon;
              return (
                <a
                  key={item.label}
                  href={item.href}
                  className="inline-flex min-h-[44px] min-w-[44px] items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-white"
                >
                  <Icon className="h-4 w-4 text-slate-500 dark:text-slate-400" />
                  {item.label}
                </a>
              );
            })}
          </nav>

          {/* Controls (Theme Switcher & Mobile Drawer Toggle) */}
          <div className="flex items-center gap-2">
            <button
              onClick={toggleTheme}
              type="button"
              aria-label={
                isDarkMode ? "Switch to Light Mode" : "Switch to Dark Mode"
              }
              className="inline-flex min-h-[44px] min-w-[44px] items-center justify-center rounded-lg border border-slate-200 bg-white p-2 text-slate-600 shadow-sm transition-colors hover:bg-slate-100 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              {isDarkMode ? (
                <Sun className="h-5 w-5 text-amber-400" />
              ) : (
                <Moon className="h-5 w-5 text-slate-600" />
              )}
            </button>

            <button
              onClick={() => setIsMobileMenuOpen((prev) => !prev)}
              type="button"
              aria-expanded={isMobileMenuOpen}
              aria-label="Toggle Navigation Menu"
              className="inline-flex min-h-[44px] min-w-[44px] items-center justify-center rounded-lg border border-slate-200 bg-white p-2 text-slate-600 transition-colors hover:bg-slate-100 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 md:hidden"
            >
              {isMobileMenuOpen ? (
                <X className="h-6 w-6" />
              ) : (
                <Menu className="h-6 w-6" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation Drawer */}
        {isMobileMenuOpen && (
          <div className="border-b border-slate-200 bg-white px-4 pt-2 pb-4 dark:border-slate-800 dark:bg-slate-900 md:hidden">
            <nav className="flex flex-col gap-1" aria-label="Mobile Navigation">
              {navItems.map((item) => {
                const Icon = item.icon;
                return (
                  <a
                    key={item.label}
                    href={item.href}
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="inline-flex min-h-[44px] items-center gap-3 rounded-lg px-3 py-2 text-base font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
                  >
                    <Icon className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                    {item.label}
                  </a>
                );
              })}
            </nav>
          </div>
        )}
      </header>

      {/* Main Area */}
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
};

// ============================================================================
// 3. MAIN APPLICATION ENTRY POINT
// ============================================================================

export function App() {
  return (
    <AppLayout>
      <div className="space-y-6">
        <header>
          <h1 className="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white sm:text-4xl">
            Tournament Dashboard
          </h1>
          <p className="mt-2 text-slate-600 dark:text-slate-400">
            Manage multi-sport youth pools, rosters, and dynamic match
            scheduling.
          </p>
        </header>

        {/* Metric Cards */}
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-800/50">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-white">
              Active Pools
            </h2>
            <p className="mt-1 text-3xl font-bold text-brand-500">12</p>
            <p className="mt-2 text-sm text-slate-500 dark:text-slate-400">
              Soccer Season 2026
            </p>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-800/50">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-white">
              Registered Players
            </h2>
            <p className="mt-1 text-3xl font-bold text-brand-500">248</p>
            <p className="mt-2 text-sm text-slate-500 dark:text-slate-400">
              Sub-registrations confirmed
            </p>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-800/50">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-white">
              Upcoming Events
            </h2>
            <p className="mt-1 text-3xl font-bold text-brand-500">4</p>
            <p className="mt-2 text-sm text-slate-500 dark:text-slate-400">
              Next: Rimouski Tournament
            </p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4">
          <Button variant="primary">Create Roster</Button>
          <Button variant="outline">Export Schedule</Button>
        </div>
      </div>
    </AppLayout>
  );
}

export default App;
