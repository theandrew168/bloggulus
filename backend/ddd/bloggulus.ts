// models

type BlogID = string;
type Blog = {
    id: BlogID;
    title: string;
    feedURL: URL;
    siteURL: URL;
};

type PostID = string;
type Post = {
    id: PostID;
    blogID: BlogID;
    url: URL;
    title: string;
    content: string;
    publishedAt: Date;
};

type TagID = string;
type Tag = {
    id: TagID;
    name: string;
};

type AccountID = string;
type Account = {
    id: AccountID;
    username: string;
    isAdmin: boolean;
    followedBlogs: BlogID[];
};

type SessionID = string;
type Session = {
    id: SessionID;
    accountID: AccountID;
    expiresAt: Date;
};

// interfaces

type FeedBlog = {
	feedURL: URL;
	siteURL: URL;
	title: string;
	posts: FeedPost[];
};

type FeedPost = {
	url: URL;
	title: string;
	content?: string;
	publishedAt: Date;
};

type FeedParser = {
    parseFeed: (url: URL, feed: string) => Promise<FeedBlog | undefined>;
};

type FeedFetcher = {
    fetchFeed: (feedURL: URL) => Promise<string | undefined>;
};

// repositories

type BlogRepository = {
    list: () => Promise<Blog[]>;
    readByID: (id: BlogID) => Promise<Blog | undefined>;
    readByFeedURL: (feedURL: URL) => Promise<Blog | undefined>;
    createOrUpdate: (blog: Blog) => Promise<void>;
    delete: (blog: Blog) => Promise<void>;
};

type PostRepository = {
    listByBlogID: (blogID: BlogID) => Promise<Post[]>;
    readByID: (id: PostID) => Promise<Post | undefined>;
    createOrUpdate: (post: Post) => Promise<void>;
    delete: (post: Post) => Promise<void>;
};

type AccountRepository = {
    list: () => Promise<Account[]>;
    readByID: (id: AccountID) => Promise<Account | undefined>;
    readByUsername: (username: string) => Promise<Account | undefined>;
    readBySessionID: (sessionID: SessionID) => Promise<Account | undefined>;
    createOrUpdate: (account: Account) => Promise<void>;
    delete: (account: Account) => Promise<void>;
};

type SessionRepository = {
    list: () => Promise<Session[]>;
    listExpired: (now: Date) => Promise<Session[]>;
    listByAccountID: (accountID: AccountID) => Promise<Session[]>;
    readByID: (id: SessionID) => Promise<Session | undefined>;
    createOrUpdate: (session: Session) => Promise<void>;
    delete: (session: Session) => Promise<void>;
};

type TagRepository = {
    list: () => Promise<Tag[]>;
    readByID: (id: TagID) => Promise<Tag | undefined>;
    createOrUpdate: (tag: Tag) => Promise<void>;
    delete: (tag: Tag) => Promise<void>;
};

// commands

async function syncAllBlogs(
    blogRepo: BlogRepository,
    postRepo: PostRepository,
    feedFetcher: FeedFetcher,
    feedParser: FeedParser,
): Promise<void> {
    const blogs = await blogRepo.list();
    await Promise.all(
        blogs.map((blog) => addOrSyncBlog(blogRepo, postRepo, feedFetcher, feedParser, blog.feedURL))
    );
}

async function addOrSyncBlog(
    blogRepo: BlogRepository,
    postRepo: PostRepository,
    feedFetcher: FeedFetcher,
    feedParser: FeedParser,
    feedURL: URL,
): Promise<void> {
    const feed = await feedFetcher.fetchFeed(feedURL);
    if (!feed) {
        // No feed content returned
        return;
    }

    const feedBlog = await feedParser.parseFeed(feedURL, feed);
    if (!feedBlog) {
        // Feed parsing failed
        return;
    }

    let existingBlog = await blogRepo.readByFeedURL(feedURL);
    if (!existingBlog) {
        existingBlog = {
            id: 'some-unique-id', // Generate a unique ID
            title: feedBlog.title,
            feedURL: feedBlog.feedURL,
            siteURL: feedBlog.siteURL,
        }
        await blogRepo.createOrUpdate(existingBlog);
    }

    const existingPosts = await postRepo.listByBlogID(existingBlog.id);
    const { postsToCreate, postsToUpdate } = comparePosts(
        existingBlog,
        existingPosts,
        feedBlog.posts,
    );

    for (const post of postsToCreate) {
        await postRepo.createOrUpdate(post);
    }
    for (const post of postsToUpdate) {
        await postRepo.createOrUpdate(post);
    }
}

function comparePosts(
    blog: Blog,
    knownPosts: Post[],
    feedPosts: FeedPost[],
): { postsToCreate: Post[]; postsToUpdate: Post[] } {
    const knownPostsByURL: Record<string, Post> = {};
    for (const post of knownPosts) {
        knownPostsByURL[post.url.toString()] = post;
    }

    const postsToCreate: Post[] = [];
    const postsToUpdate: Post[] = [];

    for (const feedPost of feedPosts) {
        const knownPost = knownPostsByURL[feedPost.url.toString()];
        if (!knownPost) {
            postsToCreate.push({
                id: 'some-unique-id', // Generate a unique ID
                blogID: blog.id,
                url: feedPost.url,
                title: feedPost.title,
                content: feedPost.content || '',
                publishedAt: feedPost.publishedAt,
            });
        } else {
            // TODO: Check if the post needs to be updated
            postsToUpdate.push(knownPost);
        }
    }

    return { postsToCreate, postsToUpdate };
}

// NOTE: This command _must_ return a session object.
async function signIn(
    accountRepo: AccountRepository,
    sessionRepo: SessionRepository,
    username: string,
): Promise<Session> {
    let account = await accountRepo.readByUsername(username);
    if (!account) {
        account = {
            id: 'some-unique-account-id', // Generate a unique account ID
            username: username,
            isAdmin: false,
            followedBlogs: [],
        }
        await accountRepo.createOrUpdate(account);
    }

    const session: Session = {
        id: 'some-unique-session-id', // Generate a unique session ID
        accountID: account.id,
        expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000), // 1 day expiration
    };
    await sessionRepo.createOrUpdate(session);

    return session;
}

async function signOut(
    sessionRepo: SessionRepository,
    sessionID: SessionID,
): Promise<void> {
    const session = await sessionRepo.readByID(sessionID);
    if (!session) {
        return;
    }

    await sessionRepo.delete(session);
}

async function deleteAccount(
    accountRepo: AccountRepository,
    accountID: AccountID,
): Promise<void> {
    const account = await accountRepo.readByID(accountID);
    if (!account) {
        throw new Error('Account not found');
    }

    await accountRepo.delete(account);
}

async function followBlog(
    accountRepo: AccountRepository,
    accountID: AccountID,
    blogID: BlogID,
): Promise<void> {
    const account = await accountRepo.readByID(accountID);
    if (!account) {
        throw new Error('Account not found');
    }

    if (!account.followedBlogs.includes(blogID)) {
        account.followedBlogs.push(blogID);
        await accountRepo.createOrUpdate(account);
    }
}

async function unfollowBlog(
    accountRepo: AccountRepository,
    accountID: AccountID,
    blogID: BlogID,
): Promise<void> {
    const account = await accountRepo.readByID(accountID);
    if (!account) {
        throw new Error('Account not found');
    }

    account.followedBlogs = account.followedBlogs.filter((id) => id !== blogID);
    await accountRepo.createOrUpdate(account);
}

async function deleteBlog(
    blogRepo: BlogRepository,
    blogID: BlogID,
): Promise<void> {
    const blog = await blogRepo.readByID(blogID);
    if (!blog) {
        throw new Error('Blog not found');
    }

    await blogRepo.delete(blog);
}

async function deletePost(
    postRepo: PostRepository,
    postID: PostID,
): Promise<void> {
    const post = await postRepo.readByID(postID);
    if (!post) {
        throw new Error('Post not found');
    }

    await postRepo.delete(post);
}

// queries

type Article = {
    title: string;
    url: URL;
    blogTitle: string;
    blogURL: URL;
    publishedAt: Date;
    tags: Tag[];
}

async function listRecentArticles(account?: Account): Promise<Article[]> {
    return [];
}

async function listRelevantArticles(search: string, account?: Account): Promise<Article[]> {
    return [];
}

async function listAccounts(accountRepo: AccountRepository): Promise<Account[]> {
    return accountRepo.list();
}

async function findAccountBySessionID(
    accountRepo: AccountRepository,
    sessionRepo: SessionRepository,
    sessionID: SessionID,
): Promise<Account | undefined> {
    const session = await sessionRepo.readByID(sessionID);
    if (!session) {
        return undefined;
    }

    return accountRepo.readByID(session.accountID);
}

async function readBlogByID(
    blogRepo: BlogRepository,
    blogID: BlogID,
): Promise<Blog | undefined> {
    return blogRepo.readByID(blogID);
}

type BlogForAccount = {
    id: BlogID;
    title: string;
    siteURL: URL;
    isFollowed: boolean;
};

async function listBlogsForAccount(
    blogRepo: BlogRepository,
    account: Account,
): Promise<BlogForAccount[]> {
    const blogs = await blogRepo.list();
    return blogs.map((blog) => ({
        id: blog.id,
        title: blog.title,
        siteURL: blog.siteURL,
        isFollowed: account.followedBlogs.includes(blog.id),
    }));
}

async function readPostByID(
    postRepo: PostRepository,
    postID: PostID,
): Promise<Post | undefined> {
    return postRepo.readByID(postID);
}