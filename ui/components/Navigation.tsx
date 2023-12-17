'use client'

import Link from "next/link"

interface Props{
    loggedin: boolean
}

function NavBar(props: Props) {
    var routes = [
        {
            path: "/challs",
            name: "Challenges",
        },
        {
            path: "/scoreboard",
            name: "Scoreboard"
        },
    ]

    if (props.loggedin){
        routes.push({
            path: "/logout",
            name: "Logout",
        });
    } else {
        routes.push({
            path: "/login",
            name: "Login",
        });
    }

    return (
            <div className="flex p-2 w-full z-[100] bg-palette-300">
                <Link href={ "/" } className="flex">
                    <img src="/static/assets/isolet.svg" className="align-center px-2"></img>
                </Link>

                <nav className="flex gap-2 justify-end relative w-full z-[100]">
                    {routes.map((item, index) => {
        
                        return (
                            <Link
                                key={item.path}
                                className="px-3 py-2 rounded-md text-sm lg:text-base font-mono relative no-underline duration-300 ease-in text-palette-200"
                                href={item.path}
                            >
                                <span>{item.name}</span>
                            </Link>
                        );
                    })}
                </nav>
            </div>
    )
}

export default NavBar;