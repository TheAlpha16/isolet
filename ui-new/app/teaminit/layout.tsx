'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { fetching, user } = useAuthStore();

	if (fetching) {
		return (
			<div className='flex justify-center w-screen animate-bounce text-2xl'>
				Loading...
			</div>
		)
	}

	if (user.userid === -1) {
		return redirect('/login');
	}

	if (user.teamid !== -1) {
		return redirect('/challenges');
	}

	return children
}