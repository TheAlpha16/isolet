import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { ThemeToggle } from "@/components/theme-toggle";
import { UserRound, LogOut, Trophy, Flag, Menu, X, LogIn, Rocket } from 'lucide-react';
import { useState } from "react";

interface Route {
path: string;
name: string;
icon: React.ReactNode;
}

function NavBar() {
	const { loggedIn, fetching, logout } = useAuthStore();
	const [ menuOpen, setMenuOpen ] = useState(false);

	const routes: Route[] = [
		{ path: "/profile", name: "Profile", icon: <UserRound size={18} /> },
		{ path: "/scoreboard", name: "Scoreboard", icon: <Trophy size={18} /> },
		{ path: "/challenges", name: "Challenges", icon: <Flag size={18} /> },
	];

	return (
		<>
		<div className="flex items-center justify-between bg-transparent p-4 font-mono border-b sm:sticky sm:top-0 z-10">
			<div className="flex items-center gap-4">
				<Link href="/">
					<div className="text-foreground text-2xl font-bold">isolet</div>
				</Link>

				{loggedIn &&
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
				}
			</div>

			<div className="items-center gap-4 hidden sm:flex">
				{!fetching && (
					<>{loggedIn ? (<>
						<ThemeToggle />
						<Button variant="ghost" size="icon" onClick={logout}>
							<LogOut size={18} />
							<span className="sr-only">Logout</span>
						</Button>
					</>) : (
					<>
						<Link href="/register">
							<Button variant="secondary">Register</Button>
						</Link>
						<Link href="/login">
							<Button>Login</Button>
						</Link>
						<ThemeToggle />
					</>
					)}</>
				)}
			</div>

			{!fetching && (
				<div className="flex sm:hidden items-center">
					<ThemeToggle />
					<Button variant="ghost" size="icon" onClick={() => setMenuOpen(!menuOpen)}>
							{menuOpen ? <X size={24} /> : <Menu size={24} />}
						<span className="sr-only">Menu</span>
					</Button>
				</div>
			)}
		</div>
		
		<div className={`bg-transparent font-mono border-b sm:hidden transition-all duration-300 ease-in-out ${
          menuOpen ? "max-h-screen opacity-100" : "max-h-0 opacity-0"}`}>
			{loggedIn && (
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
				</nav>)}
			{!loggedIn && (
				<nav className="flex flex-col gap-2 p-4">
					<Link
						className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
						href="/register"
					>
						<Rocket size={18} />
						<span>Register</span>
					</Link>
					<Link
						className="flex items-center gap-2 rounded-md p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
						href="/login"
					>
						<LogIn size={18} />
						<span>Login</span>
					</Link>
				</nav>
			)}
		</div>
		</>
	);
}

export default NavBar;