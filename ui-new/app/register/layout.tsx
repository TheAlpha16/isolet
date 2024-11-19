'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { loggedIn } = useAuthStore();

	if (loggedIn) {
		return redirect('/');
	}

	return children
}