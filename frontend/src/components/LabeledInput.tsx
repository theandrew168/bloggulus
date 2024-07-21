export type Props = {
	name: string;
	label: string;
	type?: string;
	required?: boolean;
};

export default function LabeledInput({ name, label, type, required }: Props) {
	return (
		<div>
			<label htmlFor={name} className="block text-sm font-medium">
				{label}
			</label>
			<div className="mt-2">
				<input
					id={name}
					name={name}
					placeholder={label}
					type={type}
					required={required}
					className="block rounded-md border-0 py-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-gray-800"
				/>
			</div>
		</div>
	);
}
