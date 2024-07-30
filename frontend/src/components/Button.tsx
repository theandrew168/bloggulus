import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

export type Props = {
	type?: ButtonHTMLAttributes<HTMLButtonElement>["type"];
	name?: ButtonHTMLAttributes<HTMLButtonElement>["name"];
	value?: ButtonHTMLAttributes<HTMLButtonElement>["value"];
};

export default function Button({ type, name, value, children }: PropsWithChildren<Props>) {
	return (
		<button
			type={type}
			name={name}
			value={value}
			className="rounded-md px-3 py-1.5 text-sm font-semibold bg-gray-700 text-gray-100 shadow-sm hover:bg-gray-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-700"
		>
			{children}
		</button>
	);
}
