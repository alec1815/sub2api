package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/internal/domain"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// APIKey holds the schema definition for the APIKey entity.
type APIKey struct {
	ent.Schema
}

func (APIKey) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "api_keys"},
	}
}

func (APIKey) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (APIKey) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.String("key").
			MaxLen(128).
			NotEmpty().
			Unique(),
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.Int64("group_id").
			Optional().
			Nillable(),

		// 企业分配字段（企业功能 P1 新增）
		field.Int64("assigned_to").
			Optional().
			Nillable().
			Comment("分配给企业成员的 enterprise_members.id，NULL=个人Key"),
		field.String("usage_purpose").
			MaxLen(200).
			Default("").
			Comment("用途说明"),
		field.String("bound_tool").
			MaxLen(50).
			Default("").
			Comment("绑定工具：cursor/trae/claude_code/codex/opencode/pixso/other"),
		field.String("status").
			MaxLen(20).
			Default(domain.StatusActive),
		field.Time("last_used_at").
			Optional().
			Nillable().
			Comment("Last usage time of this API key"),
		field.JSON("ip_whitelist", []string{}).
			Optional().
			Comment("Allowed IPs/CIDRs, e.g. [\"192.168.1.100\", \"10.0.0.0/8\"]"),
		field.JSON("ip_blacklist", []string{}).
			Optional().
			Comment("Blocked IPs/CIDRs"),

		// ========== Quota fields ==========
		// Quota limit in USD (0 = unlimited)
		field.Float("quota").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Quota limit in USD for this API key (0 = unlimited)"),
		// Used quota amount
		field.Float("quota_used").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Used quota amount in USD"),
		// Expiration time (nil = never expires)
		field.Time("expires_at").
			Optional().
			Nillable().
			Comment("Expiration time for this API key (null = never expires)"),

		// ========== Rate limit fields ==========
		// Rate limit configuration (0 = unlimited)
		field.Float("rate_limit_5h").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Rate limit in USD per 5 hours (0 = unlimited)"),
		field.Float("rate_limit_1d").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Rate limit in USD per day (0 = unlimited)"),
		field.Float("rate_limit_7d").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Rate limit in USD per 7 days (0 = unlimited)"),
		// Rate limit usage tracking
		field.Float("usage_5h").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Used amount in USD for the current 5h window"),
		field.Float("usage_1d").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Used amount in USD for the current 1d window"),
		field.Float("usage_7d").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Used amount in USD for the current 7d window"),
		// Window start times
		field.Time("window_5h_start").
			Optional().
			Nillable().
			Comment("Start time of the current 5h rate limit window"),
		field.Time("window_1d_start").
			Optional().
			Nillable().
			Comment("Start time of the current 1d rate limit window"),
		field.Time("window_7d_start").
			Optional().
			Nillable().
			Comment("Start time of the current 7d rate limit window"),
	}
}

func (APIKey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("api_keys").
			Field("user_id").
			Unique().
			Required(),
		edge.From("group", Group.Type).
			Ref("api_keys").
			Field("group_id").
			Unique(),
		edge.To("usage_logs", UsageLog.Type),

		// 企业 Key 分配关系 —— assigned_to → enterprise_members.id
		edge.From("assigned_member", EnterpriseMember.Type).
			Ref("assigned_keys").
			Field("assigned_to").
			Unique().
			Comment("Key 被分配给的企业成员"),

		// M:N 分组关联（企业功能 P1 新增，替代 group_id 1:1）
		edge.To("key_groups", Group.Type).
			Through("api_key_groups", APIKeyGroup.Type).
			Comment("Key 关联的多个分组（M:N）"),
	}
}

func (APIKey) Indexes() []ent.Index {
	return []ent.Index{
		// key 字段已在 Fields() 中声明 Unique()，无需重复索引
		index.Fields("user_id"),
		index.Fields("group_id"),
		index.Fields("status"),
		index.Fields("deleted_at"),
		index.Fields("last_used_at"),
		// Index for quota queries
		index.Fields("quota", "quota_used"),
		index.Fields("expires_at"),
		// 企业 Key 查询索引（企业功能 P1 新增）
		index.Fields("assigned_to"),
		index.Fields("assigned_to", "status"),
	}
}
