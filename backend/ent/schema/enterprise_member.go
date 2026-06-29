package schema

import (
	"time"

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

// EnterpriseMember holds the schema definition for the EnterpriseMember entity.
// 本期约束：一人一企业（1:1），一个 user_id 最多一条 active 记录。
type EnterpriseMember struct {
	ent.Schema
}

func (EnterpriseMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_members"},
	}
}

func (EnterpriseMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (EnterpriseMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("enterprise_id").
			Comment("所属企业 ID"),
		field.Int64("user_id").
			Comment("关联用户 ID"),

		// 企业内角色（详见 02-核心设计决策.md 决策 2：两层角色体系）
		field.String("role").
			MaxLen(20).
			Default(domain.EnterpriseRoleMember).
			Comment("企业内角色：enterprise_admin / enterprise_member"),

		// 成员状态
		field.String("status").
			MaxLen(20).
			Default(domain.StatusActive).
			Comment("active / unbound"),

		// 所属部门（可选）
		field.Int64("department_id").
			Optional().
			Nillable().
			Comment("所属部门 ID"),

		// 成员级限制
		field.Int("concurrency").
			Default(0).
			Comment("成员级并发上限，0=不限制"),
		field.Int("rpm_limit").
			Default(0).
			Comment("成员级 RPM 上限，0=不限制"),

		// 备注
		field.String("notes").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default("").
			Comment("备注"),

		// 加入/解绑时间
		field.Time("joined_at").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("加入企业时间"),
		field.Time("unbound_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("解绑时间"),
	}
}

func (EnterpriseMember) Edges() []ent.Edge {
	return []ent.Edge{
		// 所属企业
		edge.From("enterprise", Enterprise.Type).
			Ref("members").
			Field("enterprise_id").
			Required().
			Unique(),
		// 关联用户
		edge.From("user", User.Type).
			Ref("enterprise_members").
			Field("user_id").
			Required().
			Unique(),
		// 所属部门
		edge.From("department", Department.Type).
			Ref("members").
			Field("department_id").
			Unique(),
		// 分配给该成员的 API Key（assigned_to → enterprise_members.id）
		edge.To("assigned_keys", APIKey.Type).
			Comment("分配给该成员的 API Key"),
	}
}

func (EnterpriseMember) Indexes() []ent.Index {
	return []ent.Index{
		// (enterprise_id, user_id) 唯一（WHERE deleted_at IS NULL）—— 防同一企业重复加入
		// 部分唯一索引由迁移 SQL 实现
		index.Fields("enterprise_id", "user_id").
			Unique(),
		// user_id 部分唯一索引（WHERE status='active' AND deleted_at IS NULL）—— 本期 1:1 约束
		// 由迁移 SQL 实现: CREATE UNIQUE INDEX ON enterprise_members(user_id) WHERE status='active' AND deleted_at IS NULL
		index.Fields("user_id"),
		index.Fields("enterprise_id"),
		index.Fields("role"),
		index.Fields("enterprise_id", "role"),
		index.Fields("department_id"),
		index.Fields("deleted_at"),
	}
}
