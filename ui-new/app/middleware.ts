import { useAuthStore } from '@/store/authStore';
import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
    const { loggedIn, user, logout } = useAuthStore();

    if (!loggedIn) {
        return NextResponse.redirect(new URL("/login", request.url));
    }

    if (!user) {
        logout();
        return NextResponse.redirect(new URL("/login", request.url));
    }

    if (request.nextUrl.pathname.startsWith("/onboarding")) {
        if (user?.teamid !== -1) {
            return NextResponse.redirect(new URL("/", request.url));
        }
    }

    if (user?.teamid === -1) {
        return NextResponse.redirect(new URL("/onboarding", request.url));
    }

    return NextResponse.next();
}

