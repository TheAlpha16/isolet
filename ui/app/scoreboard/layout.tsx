'use client'

import React from 'react';
import { useAuthStore } from '@/store/authStore';
import { redirect } from 'next/navigation';
import { ScoreboardSkeleton } from '@/components/skeletons/scoreboard';

export default function RootLayout({ children, }: { children: React.ReactNode }) {
	const { user, fetching } = useAuthStore();

	if (fetching) {
		return <ScoreboardSkeleton />;
	}

	if (user.userid === -1) {
		return redirect('/login');
	}

	if (user.teamid === -1) {
		return redirect('/teaminit');
	}

	return children
}
