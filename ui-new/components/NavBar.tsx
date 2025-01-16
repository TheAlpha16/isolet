import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { ThemeToggle } from "@/components/theme-toggle";
import { CircleUserRound, LogOut, Users, Users2, Trophy, Flag } from 'lucide-react';

interface Route {
path: string;
name: string;
icon: React.ReactNode;
}

function NavBar() {
	const { loggedIn, fetching, user, logout } = useAuthStore();

	const routes: Route[] = [
		{ path: "/teams", name: "Teams", icon: <Users2 size={18} /> },
		{ path: "/scoreboard", name: "Scoreboard", icon: <Trophy size={18} /> },
		{ path: "/challenges", name: "Challenges", icon: <Flag size={18} /> },
	];

	return (
		<div className="flex items-center justify-between bg-transparent p-4 font-mono border-b sticky top-0 z-10">
			<div className="flex items-center gap-4">
				<Link href="/">
					<div className="text-foreground text-2xl font-bold">isolet</div>
				</Link>

				{loggedIn &&
					<nav className="flex items-center gap-2">
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

			<div className="flex items-center gap-4">
				{!fetching && (
					<>{loggedIn ? (<>
						<Link href="/profile">
							<Button variant="ghost" size="icon">
								<CircleUserRound size={18} />
								<span className="sr-only">Profile</span>
							</Button>
						</Link>
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
					</>
					)}</>
				)}
				<ThemeToggle />
			</div>
		</div>
	);
}

export default NavBar;