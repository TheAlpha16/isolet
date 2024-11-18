import { useAuthStore } from "@/store/authStore";
import FormButton from "@/components/extras/buttons";
import Link from "next/link";

interface Route {
	path: string;
	name: string;
}

function NavBar() {
	const { loggedIn, user, logout } = useAuthStore();
	const routes: Route[] = [
		{ path: "/users", name: "Users" },
		{ path: "/teams", name: "Teams" },
		{ path: "/scoreboard", name: "Scoreboard" },
		{ path: "/challenges", name: "Challenges" },
	];

	return (
		<div className="flex gap-4 bg-transparent p-4 font-mono items-center">
			<Link href="/">
				<div className="text-white text-2xl font-bold">isolet</div>
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

			<div className="ml-auto">
				{loggedIn ? (
					<div> profile </div>
				) : (
					<div className="flex gap-2">
						<FormButton type="button" variant="secondary">
							<Link href="/register">Register</Link>
						</FormButton>
						<FormButton type="button" variant="primary">
							<Link href="/login">Login</Link>
						</FormButton>
					</div>
				)}
			</div>
		</div>
	);
}

export default NavBar;
