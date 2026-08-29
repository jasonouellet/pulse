import React, { useState, useEffect } from "react";
import { DropdownMenu } from "radix-ui";
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
  ChevronDown,
  Check,
  Building2,
  UsersRound,
} from "lucide-react";
import { ROLE_LABELS, useSession, type UserRole } from "../../context/session";

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

  const {
    userName,
    userInitials,
    grants,
    clubs,
    activeGrant,
    setActiveRole,
    setActiveClub,
  } = useSession();

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

  const rolesForActiveClub = grants
    .filter((g) => g.club.id === activeGrant.club.id)
    .map((g) => g.role);

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

          {/* Right Actions (Theme Switcher + Avatar Menu + Mobile Toggle) */}
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

            {/* Avatar + role/club submenu — the ONLY place club/role switching
                happens. The name/subtitle next to the avatar is passive
                display text, not a separate control. */}
            <DropdownMenu.Root>
              <DropdownMenu.Trigger asChild>
                <button
                  type="button"
                  className="inline-flex min-h-[44px] items-center gap-2 rounded-lg px-2 py-1.5 hover:bg-slate-100 focus:outline-none focus:ring-2 focus:ring-brand-500 dark:hover:bg-slate-800"
                >
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-brand-100 text-sm font-medium text-brand-700 dark:bg-brand-500/20 dark:text-brand-300">
                    {userInitials}
                  </div>
                  <div className="hidden text-left sm:block">
                    <p className="text-sm font-medium leading-tight text-slate-900 dark:text-white">
                      {userName}
                    </p>
                    <p className="text-xs leading-tight text-slate-500 dark:text-slate-400">
                      {ROLE_LABELS[activeGrant.role]} · {activeGrant.club.name}
                    </p>
                  </div>
                  <ChevronDown
                    className="hidden h-4 w-4 text-slate-400 sm:block"
                    aria-hidden="true"
                  />
                </button>
              </DropdownMenu.Trigger>

              <DropdownMenu.Portal>
                <DropdownMenu.Content
                  align="end"
                  sideOffset={8}
                  className="z-50 w-64 rounded-xl border border-slate-200 bg-white p-1.5 shadow-lg dark:border-slate-700 dark:bg-slate-800"
                >
                  <DropdownMenu.Label className="px-2 py-1.5 text-xs font-medium text-slate-400 dark:text-slate-500">
                    Rôle actif · {activeGrant.club.name}
                  </DropdownMenu.Label>
                  {rolesForActiveClub.map((role: UserRole) => (
                    <DropdownMenu.Item
                      key={role}
                      onSelect={() => setActiveRole(role)}
                      className="flex cursor-pointer items-center gap-2 rounded-lg px-2 py-2 text-sm text-slate-700 outline-none data-[highlighted]:bg-slate-100 dark:text-slate-200 dark:data-[highlighted]:bg-slate-700"
                    >
                      <UsersRound
                        className="h-4 w-4 text-slate-400"
                        aria-hidden="true"
                      />
                      <span className="flex-1">{ROLE_LABELS[role]}</span>
                      {role === activeGrant.role && (
                        <Check
                          className="h-4 w-4 text-brand-500"
                          aria-hidden="true"
                        />
                      )}
                    </DropdownMenu.Item>
                  ))}

                  <DropdownMenu.Separator className="my-1.5 h-px bg-slate-200 dark:bg-slate-700" />

                  <DropdownMenu.Label className="px-2 py-1.5 text-xs font-medium text-slate-400 dark:text-slate-500">
                    Club actif
                  </DropdownMenu.Label>
                  {clubs.map((club) => (
                    <DropdownMenu.Item
                      key={club.id}
                      onSelect={() => setActiveClub(club.id)}
                      className="flex cursor-pointer items-center gap-2 rounded-lg px-2 py-2 text-sm text-slate-700 outline-none data-[highlighted]:bg-slate-100 dark:text-slate-200 dark:data-[highlighted]:bg-slate-700"
                    >
                      <Building2
                        className="h-4 w-4 text-slate-400"
                        aria-hidden="true"
                      />
                      <span className="flex-1">{club.name}</span>
                      {club.id === activeGrant.club.id && (
                        <Check
                          className="h-4 w-4 text-brand-500"
                          aria-hidden="true"
                        />
                      )}
                    </DropdownMenu.Item>
                  ))}

                  <DropdownMenu.Separator className="my-1.5 h-px bg-slate-200 dark:bg-slate-700" />

                  <DropdownMenu.Item className="flex cursor-pointer items-center gap-2 rounded-lg px-2 py-2 text-sm text-slate-700 outline-none data-[highlighted]:bg-slate-100 dark:text-slate-200 dark:data-[highlighted]:bg-slate-700">
                    <Settings
                      className="h-4 w-4 text-slate-400"
                      aria-hidden="true"
                    />
                    Paramètres du compte
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Portal>
            </DropdownMenu.Root>

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
