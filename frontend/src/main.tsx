import React from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import SiteLayout from "./layouts/SiteLayout";
import ErrorPage from "./pages/ErrorPage";
import IndexPage, { indexPageLoader } from "./pages/IndexPage";
import LoginPage, { loginPageAction } from "./pages/LoginPage";
import RegisterPage, { registerPageAction } from "./pages/RegisterPage";
import BlogsPage, { blogsPageAction, blogsPageLoader } from "./pages/BlogsPage";
import TagsPage, { tagsPageAction, tagsPageLoader } from "./pages/TagsPage";
import BlogPage, { blogPageAction, blogPageLoader } from "./pages/BlogPage";
import PostPage, { postPageAction, postPageLoader } from "./pages/PostPage";

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
			{
				path: "/blogs/:blogID",
				loader: blogPageLoader,
				action: blogPageAction,
				element: <BlogPage />,
			},
			{
				path: "/blogs/:blogID/posts/:postID",
				loader: postPageLoader,
				action: postPageAction,
				element: <PostPage />,
			},
			{
				path: "/tags",
				loader: tagsPageLoader,
				action: tagsPageAction,
				element: <TagsPage />,
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
