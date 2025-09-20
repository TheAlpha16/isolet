"use client";

import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { useEventStore, useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { Flag, LogIn, LogOut, Menu, Rocket, Trophy, UserRound, X } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

interface Route {
  path: string;
  name: string;
  icon: React.ReactNode;
}

function NavBar() {
  const { user, meLoading } = useProfileStore();
  const { logout } = useProfileStore();
  const {
    info: { name },
  } = useEventStore();
  const [menuOpen, setMenuOpen] = useState(false);

  const routes: Route[] = [
    { path: UI_ROUTES.profile.home, name: "Profile", icon: <UserRound size={18} /> },
    { path: UI_ROUTES.score.home, name: "Scoreboard", icon: <Trophy size={18} /> },
    { path: UI_ROUTES.challenge.home, name: "Challenges", icon: <Flag size={18} /> },
  ];

  return (
    <>
      <div className="flex items-center justify-between bg-background p-4 font-mono border-b sm:sticky sm:top-0 z-10">
        <div className="flex items-center gap-4 min-h-[40px]">
          <Link href={UI_ROUTES.home} className="flex items-center gap-2">
            <div className="text-foreground text-2xl font-bold">{name}</div>
          </Link>

          {user && (
            <nav className="sm:flex items-center gap-2 hidden">
              {routes.map(({ path, name, icon }) => (
                <Link
                  key={path}
                  className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
                  href={path}
                >
                  {icon}
                  <span>{name}</span>
                </Link>
              ))}
            </nav>
          )}
        </div>

        <div className="items-center gap-4 hidden sm:flex">
          {!meLoading && (
            <>
              {user ? (
                <>
                  <ThemeToggle />
                  <Button variant="ghost" size="icon" onClick={logout}>
                    <LogOut size={18} />
                    <span className="sr-only">Logout</span>
                  </Button>
                </>
              ) : (
                <>
                  <Link href={UI_ROUTES.onboard.user}>
                    <Button variant="secondary">Register</Button>
                  </Link>
                  <Link href={UI_ROUTES.auth.login}>
                    <Button>Login</Button>
                  </Link>
                  <ThemeToggle />
                </>
              )}
            </>
          )}
        </div>

        {!meLoading && (
          <div className="flex sm:hidden items-center">
            <ThemeToggle />
            <Button variant="ghost" size="icon" onClick={() => setMenuOpen(!menuOpen)}>
              {menuOpen ? <X size={24} /> : <Menu size={24} />}
              <span className="sr-only">Menu</span>
            </Button>
          </div>
        )}
      </div>

      {menuOpen && (
        <div
          className={`bg-transparent font-mono border-b sm:hidden transition-all duration-300 ease-in-out `}
        >
          {user && (
            <nav className="flex flex-col gap-2 p-4">
              {routes.map(({ path, name, icon }) => (
                <Link
                  key={path}
                  className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
                  href={path}
                >
                  {icon}
                  <span>{name}</span>
                </Link>
              ))}
              <Link
                className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
                href="#"
                onClick={logout}
              >
                <LogOut size={18} />
                <span>Logout</span>
              </Link>
            </nav>
          )}
          {!user && (
            <nav className="flex flex-col gap-2 p-4">
              <Link
                className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
                href={UI_ROUTES.onboard.user}
              >
                <Rocket size={18} />
                <span>Register</span>
              </Link>
              <Link
                className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
                href={UI_ROUTES.auth.login}
              >
                <LogIn size={18} />
                <span>Login</span>
              </Link>
            </nav>
          )}
        </div>
      )}
    </>
  );
}

export default NavBar;
