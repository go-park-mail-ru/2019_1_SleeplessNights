package query

const CreateNewForum = `
INSERT INTO "Forums" (slug, title, "user-id")
VALUES ($1, $2, $3)
`