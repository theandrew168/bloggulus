import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./pages/ErrorPage";
import IndexPage, { indexPageLoader } from "./pages/IndexPage";
import SiteLayout from "./layouts/SiteLayout";
import LoginPage, { loginPageAction } from "./pages/LoginPage";
import RegisterPage, { registerPageAction } from "./pages/RegisterPage";
import BlogsPage, { blogsPageAction, blogsPageLoader } from "./pages/BlogsPage";

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
				action: loginPageAction,
				element: <LoginPage />,
			},
			{
				path: "/register",
				action: registerPageAction,
				element: <RegisterPage />,
			},
			{
				path: "/blogs",
				loader: blogsPageLoader,
				action: blogsPageAction,
				element: <BlogsPage />,
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
