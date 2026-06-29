package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EnterpriseSubscription holds the schema definition for the EnterpriseSubscription entity.
// 企业套餐表，与企业管理员个人套餐完全隔离。
// 无 SoftDeleteMixin——通过 status 字段管理生命周期（active/expired/suspended）。
type EnterpriseSubscription struct {
	ent.Schema
}

func (EnterpriseSubscription) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_subscriptions"},
	}
}

func (EnterpriseSubscription) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (EnterpriseSubscription) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("enterprise_id").
			Comment("所属企业 ID"),
		field.Int64("group_id").
			Comment("分组 ID"),
		field.Int64("plan_id").
			Comment("套餐计划 ID"),

		// 有效期
		field.Time("starts_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("套餐生效时间"),
		field.Time("expires_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("套餐过期时间"),

		// 状态
		field.String("status").
			MaxLen(20).
			Default("active").
			Comment("active / expired / suspended"),

		// 用量统计（企业套餐累计使用量）
		field.Float("daily_usage_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0).
			Comment("当日用量（USD）"),
		field.Float("weekly_usage_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0).
			Comment("当周用量（USD）"),
		field.Float("monthly_usage_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0).
			Comment("当月用量（USD）"),
	}
}

func (EnterpriseSubscription) Edges() []ent.Edge {
	return []ent.Edge{
		// 所属企业
		edge.From("enterprise", Enterprise.Type).
			Ref("subscriptions").
			Field("enterprise_id").
			Required().
			Unique(),
		// 关联分组
		edge.To("group", Group.Type).
			Unique().
			Required().
			Field("group_id"),
		// 关联套餐计划
		edge.To("plan", SubscriptionPlan.Type).
			Unique().
			Required().
			Field("plan_id"),
	}
}

func (EnterpriseSubscription) Indexes() []ent.Index {
	return []ent.Index{
		// 查企业有效套餐
		index.Fields("enterprise_id", "status"),
		// 查企业某分组套餐
		index.Fields("enterprise_id", "group_id"),
	}
}
