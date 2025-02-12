import React from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { Lightbulb } from 'lucide-react';

export const HintToastContainer = () => {
	return (
		<ToastContainer
			containerId="hint-toast"
			position="bottom-right"
			autoClose={5000}
			hideProgressBar={true}
			newestOnTop={false}
			closeOnClick
			rtl={false}
			pauseOnFocusLoss
			draggable
			pauseOnHover
		/>
	);
};

export function showHint(hint: string) {
	toast.info(hint, {
		containerId: "hint-toast",
		icon: <Lightbulb className="text-green-500" />,
		className: 'hint-toast',
	});
};