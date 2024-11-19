import React from "react";

type HighLightButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
	variant?: "primary" | "secondary";
	className?: string; 
};

export default function FormButton({
	children,
	variant = "primary",
	className = "",
	...props
}: HighLightButtonProps) {
    const baseStyles =
		"rounded-md border border-solid transition-colors flex items-center justify-center text-sm sm:text-base h-10 px-4 sm:px-5";
    
	const variantStyles =
		variant === "primary"
			? "border-transparent bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc]"
			: "border-black/[.08] dark:border-white/[.145] hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent sm:min-w-44";

	return (
		<button
			className={`${baseStyles} ${variantStyles} ${className}`}
			{...props}
		>
			{children}
		</button>
	);
}
