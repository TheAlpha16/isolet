'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { loggedIn, user, fetching } = useAuthStore();

	if (fetching) {
		return <div>Loading...</div>
	}

	if (!loggedIn) {
		return redirect('/login');
	}

	if (user?.teamid === -1) {
		return redirect('/teaminit');
	}

	return children
}