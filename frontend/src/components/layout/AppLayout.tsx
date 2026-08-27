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

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
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
      {/* Header / Navbar */}
      <header className="sticky top-0 z-40 w-full border-b border-slate-200 bg-white/80 backdrop-blur-md dark:border-slate-800 dark:bg-slate-900/80">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          {/* Brand Logo */}
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-500 font-bold text-white shadow-md shadow-brand-500/20">
              P
            </div>
            <span className="text-xl font-bold tracking-tight text-slate-900 dark:text-white">
              PULSE<span className="text-brand-500">.os</span>
            </span>
          </div>

          {/* Desktop Navigation */}
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

          {/* Right Actions (Theme Switcher + Mobile Toggle) */}
          <div className="flex items-center gap-2">
            {/* Theme Toggle Button */}
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

            {/* Mobile Menu Button */}
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

      {/* Main Content Area */}
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
};
