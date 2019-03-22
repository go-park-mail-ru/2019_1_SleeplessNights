package query

const findForumBySlugQuery = `
SELECT *
FROM Forums f
WHERE f.slug = $1
`

type FindForumBySlug struct {
	Slug string
}

func (sql *FindForumBySlug)GetQuery() string {
	return findForumBySlugQuery
}

func (sql *FindForumBySlug)GetArgs() []interface{} {
	return []interface{} {sql.Slug}
}
