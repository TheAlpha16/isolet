import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/extras/buttons";
import Link from "next/link";
import Image from "next/image";
import { ThemeToggle } from "@/components/theme-toggle";

interface Route {
	path: string;
	name: string;
}

function NavBar() {
	const { loggedIn, fetching, user, logout } = useAuthStore();
	const routes: Route[] = [
		{ path: "/users", name: "Users" },
		{ path: "/teams", name: "Teams" },
		{ path: "/scoreboard", name: "Scoreboard" },
		{ path: "/challenges", name: "Challenges" },
	];

	return (
		<div className="flex gap-4 bg-transparent p-4 font-mono items-center">
			<Link href="/">
				<div className="text-foreground text-2xl font-bold">isolet</div>
			</Link>

			{loggedIn && (
				<nav className="flex h-full items-center gap-2">
					{routes.map(({ path, name }) => (
						<Link
							key={path}
							className="hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent p-2 rounded-md transition-colors"
							href={path}
						>
							<span>{name}</span>
						</Link>
					))}
				</nav>
			)}

			{!fetching && (
				<div className="ml-auto">
					{loggedIn ? (
						<div className="flex gap-4">
							<Link href="/profile">
								<Image
									className="svg-icon"
									src="/profile.svg"
									alt="profile"
									width={28}
									height={28}
								></Image>
							</Link>
							<Image
								className="svg-icon hover:cursor-pointer"
								src="/logout.svg"
								alt="logout"
								width={28}
								height={28}
								onClick={logout}
							></Image>
						</div>
					) : (
						<div className="flex gap-2">
							<Link href="/register">
								<Button variant={"secondary"}>
									Register
								</Button>
							</Link>
							<Link href="/login">
								<Button>
									Login
								</Button>
							</Link>
						</div>
					)}
				</div>
			)}
			<ThemeToggle />
		</div>
	);
}

export default NavBar;
