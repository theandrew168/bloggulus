import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./pages/ErrorPage";
import IndexPage, { loader } from "./pages/IndexPage";

const router = createBrowserRouter([
	{
		path: "/",
		loader: loader,
		element: <IndexPage />,
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
