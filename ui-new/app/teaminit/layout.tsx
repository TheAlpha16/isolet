'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';
import { FormSkeleton } from '@/components/skeletons/form';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { fetching, user } = useAuthStore();

	if (fetching) {
		return <FormSkeleton />
	}

	if (user.userid === -1) {
		return redirect('/login');
	}

	if (user.teamid !== -1) {
		return redirect('/challenges');
	}

	return children
}