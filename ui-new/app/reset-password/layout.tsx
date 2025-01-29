'use client'

import React, { Suspense } from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { user } = useAuthStore();

	if (user.userid !== -1) {
		return redirect('/');
	}

	return (
		<Suspense fallback={<div>Loading...</div>}>
			{children}
		</Suspense>
	)
}