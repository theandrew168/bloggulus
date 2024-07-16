import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./ErrorPage";
import Root, { loader } from "./Root";

const router = createBrowserRouter([
	{
		path: "/",
		loader: loader,
		element: <Root />,
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
