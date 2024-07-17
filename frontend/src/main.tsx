import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./pages/ErrorPage";
import IndexPage, { indexPageLoader } from "./pages/IndexPage";
import SiteLayout from "./layouts/SiteLayout";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";

const router = createBrowserRouter([
	{
		element: <SiteLayout />,
		errorElement: <ErrorPage />,
		children: [
			{
				path: "/",
				loader: indexPageLoader,
				element: <IndexPage />,
			},
			{
				path: "/login",
				element: <LoginPage />,
			},
			{
				path: "/register",
				element: <RegisterPage />,
			},
		],
	},
]);

const container = document.getElementById("app")!;
const app = createRoot(container);
app.render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>,
);
