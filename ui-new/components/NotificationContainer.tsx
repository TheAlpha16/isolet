import React from 'react';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export const NotificationContainer = () => {
	return (
		<ToastContainer
			containerId="notification-toast"
			position="top-right"
			toastClassName="notification-toast"
		/>
	);
};
