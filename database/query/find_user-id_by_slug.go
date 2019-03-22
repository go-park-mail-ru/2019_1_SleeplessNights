package query

const findUserIdBySlugQuery = `
SELECT id
FROM Users u
WHERE u.slug = $1
`

type FindUserIdBySlug struct {
	Slug string
}

func (sql *FindUserIdBySlug)GetQuery() string {
	return findUserIdBySlugQuery
}

func (sql *FindUserIdBySlug)GetArgs() []interface{} {
	return []interface{} {sql.Slug}
}