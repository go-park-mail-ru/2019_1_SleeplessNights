package query

const FindUserIdByNicknameQuery = `
SELECT id
FROM "Users" u
WHERE u.nickname = $1
`

