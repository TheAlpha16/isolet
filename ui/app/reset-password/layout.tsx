'use client'

import React, { Suspense } from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';
import { Skeleton } from '@/components/ui/skeleton';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { user } = useAuthStore();

	if (user.userid !== -1) {
		return redirect('/');
	}

	return (
		<Suspense fallback={
			<div className="w-full h-full flex items-center justify-center">
				<Skeleton className="w-[350px] h-[350px]" />
			</div>
		}>
			{children}
		</Suspense>
	)
}