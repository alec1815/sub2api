package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// APIKeyGroup holds the edge schema definition for the api_key_groups M:N relationship.
// 替换 api_keys.group_id 的 1:1 关联，支持一个 Key 关联多个分组。
type APIKeyGroup struct {
	ent.Schema
}

func (APIKeyGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "api_key_groups"},
	}
}

func (APIKeyGroup) Fields() []ent.Field {
	return []ent.Field{
		// id 字段由 Ent 自动生成（BIGSERIAL PRIMARY KEY），无需手动定义
		field.Int64("api_key_id").
			Comment("API Key ID"),
		field.Int64("group_id").
			Comment("分组 ID"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("创建时间"),
	}
}

func (APIKeyGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("api_key", APIKey.Type).
			Unique().
			Required().
			Field("api_key_id"),
		edge.To("group", Group.Type).
			Unique().
			Required().
			Field("group_id"),
	}
}

func (APIKeyGroup) Indexes() []ent.Index {
	return []ent.Index{
		// (api_key_id, group_id) — 唯一，防重复关联
		index.Fields("api_key_id", "group_id").
			Unique(),
		// 查某分组下所有 Key
		index.Fields("group_id"),
	}
}
