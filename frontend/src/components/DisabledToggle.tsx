import { Switch } from "@headlessui/react";

export default function DisabledToggle() {
	return (
		<Switch
			disabled
			className="group relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent bg-gray-200 transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-gray-600 focus:ring-offset-2 data-[checked]:bg-gray-600"
		>
			<span className="sr-only">Use setting</span>
			<span
				aria-hidden="true"
				className="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out group-data-[checked]:translate-x-5"
			/>
		</Switch>
	);
}
