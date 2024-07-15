import React from "react";
import { createRoot } from "react-dom/client";

const app = document.getElementById("app")!;
const root = createRoot(app);
root.render(
	<React.StrictMode>
		<h1 className="sans-serif">Hello, world</h1>
	</React.StrictMode>,
);
