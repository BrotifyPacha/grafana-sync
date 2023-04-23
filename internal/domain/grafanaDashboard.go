package domain

type GrafanaDashboard struct {
	Uid      string
	Title    string
	FolderId int
}

const (
	TARGET_SQL    = "sql"
	TARGET_PROMQL = "promql"
)

type GrafanaDashboardDetails struct {
	Title   string
	Panels  []GrafanaPanel
	RawData []byte
}

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
	if t.Expr != "" {
		result.Type = TARGET_PROMQL
		result.Expression = t.Expr
		return result, true
	}
	return result, false
}
