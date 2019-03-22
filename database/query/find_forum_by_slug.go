package query

const FindForumBySlugQuery = `
SELECT f.posts    AS PostsCount,
       f.slug     AS ForumSlug,
       f.threads  AS ThreadsCount,
       f.title    AS ForumTitle,
	   u.nickname AS UserNickname
FROM "Forums" f
JOIN "Users" u ON f."user-id" = u.id
WHERE f.slug = $1
`
