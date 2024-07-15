import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./ErrorPage";

const router = createBrowserRouter([
	{
		path: "/",
		element: <h1 className="sans-serif">Hello, world</h1>,
		errorElement: <ErrorPage />,
	},
]);

const container = document.getElementById("root")!;
const root = createRoot(container);
root.render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>,
);
