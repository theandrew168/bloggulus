import type { PropsWithChildren } from "react";
import { Link } from "react-router-dom";
import Button from "./Button";

export type Props = {
	to: string;
};

export default function ButtonLink({ to, children }: PropsWithChildren<Props>) {
	return (
		<Link to={to}>
			<Button>{children}</Button>
		</Link>
	);
}
