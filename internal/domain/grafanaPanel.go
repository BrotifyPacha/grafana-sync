package domain

type GrafanaPanel struct {
	Uid     string
	Title   string
	Targets []Target
}

type Target struct {
	RefId    string
	Expr     string
	RawSql   string
	RawQuery interface{}
}

const (
	TARGET_SQL    = "sql"
)

type Query struct {
	Title      string
	Type       string
	Expression string
}

func (p *GrafanaPanel) GetQueries() []Query {
	result := []Query{}
	for _, target := range p.Targets {
		if query, ok := target.GetQuery(); ok {
			result = append(result, query)
		}
	}
	return result
}

func (t *Target) GetQuery() (result Query, ok bool) {
	result.Title = t.RefId
	rawQueryBool, ok := t.RawQuery.(bool)
	if ok && rawQueryBool && t.RawSql != "" {
		result.Type = TARGET_SQL
		result.Expression = t.RawSql
		return result, true
	}
	return result, false
}
