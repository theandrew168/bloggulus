import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

export type Props = {
	type?: ButtonHTMLAttributes<HTMLButtonElement>["type"];
};

export default function Button({ type, children }: PropsWithChildren<Props>) {
	return (
		<button
			type={type}
			className="rounded-md px-3 py-1.5 text-sm font-semibold bg-gray-700 text-gray-100 shadow-sm hover:bg-gray-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-700"
		>
			{children}
		</button>
	);
}
