'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';
import { FormSkeleton } from '@/components/skeletons/form';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { user, fetching } = useAuthStore();

	if (fetching) {
		return <FormSkeleton />
	}

	if (user.userid !== -1) {
		return redirect('/');
	}

	return children
}