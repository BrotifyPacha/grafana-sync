package domain

import (
	"reflect"
	"testing"
)

const (
	sql_query    = "select * from table_name"
	promql_query = `metric_name{ label="A" }`
)

func TestGrafanaPanel_GetQueries(t *testing.T) {
	type fields struct {
		Uid     string
		Title   string
		Targets []Target
	}
	tests := []struct {
		name   string
		fields fields
		want   []Query
	}{
		{
			name: "panel with queries of multiple types",
			fields: fields{
				Uid:   "uid",
				Title: "Panel #1",
				Targets: []Target{
					{
						RefId:    "A",
						Expr:     "",
						RawSql:   sql_query,
						RawQuery: true,
					},
					{
						RefId:    "B",
						Expr:     promql_query,
						RawSql:   "",
						RawQuery: nil,
					},
					{
						RefId:    "C",
						Expr:     "",
						RawSql:   "",
						RawQuery: false,
					},
				},
			},
			want: []Query{
				{
					Title:      "A",
					Type:       "sql",
					Expression: sql_query,
				},
				{
					Title:      "B",
					Type:       "promql",
					Expression: promql_query,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GrafanaPanel{
				Uid:     tt.fields.Uid,
				Title:   tt.fields.Title,
				Targets: tt.fields.Targets,
			}
			if got := p.GetQueries(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrafanaPanel.GetQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}
