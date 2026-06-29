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

// Enterprise holds the schema definition for the Enterprise entity.
type Enterprise struct {
	ent.Schema
}

func (Enterprise) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprises"},
	}
}

func (Enterprise) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (Enterprise) Fields() []ent.Field {
	return []ent.Field{
		// 基本信息
		field.String("name").
			MaxLen(255).
			NotEmpty().
			Comment("企业名称"),
		field.String("short_name").
			MaxLen(100).
			Default("").
			Comment("企业简称"),
		field.String("credit_code").
			MaxLen(50).
			Default("").
			Comment("统一社会信用代码"),
		field.String("address").
			MaxLen(500).
			Default("").
			Comment("企业地址"),
		field.String("scale").
			MaxLen(20).
			Default("").
			Comment("企业规模：micro/small/medium/large"),
		field.String("industry").
			MaxLen(50).
			Default("").
			Comment("所属行业"),

		// 企业层级
		// 注意：parent_id 直接作为字段使用，不定义 edge
		// 因为 parent_id=0 表示顶级企业，无对应记录，不适合作为 FK
		field.Int64("parent_id").
			Default(0).
			Comment("父企业 ID，0=顶级企业"),

		// 状态控制
		field.String("status").
			MaxLen(20).
			Default(domain.StatusActive).
			Comment("active / disabled"),

		// 联系人信息
		field.String("contact_name").
			MaxLen(100).
			Default("").
			Comment("联系人姓名"),
		field.String("contact_phone").
			MaxLen(50).
			Default("").
			Comment("联系人手机号"),
		field.String("contact_email").
			MaxLen(255).
			Default("").
			Comment("联系人邮箱"),

		// 备注
		field.String("notes").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default("").
			Comment("备注"),

		// 企业资金池（独立于管理员个人余额）
		field.Float("balance").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("企业余额（企业独立资金池）"),
		field.Float("total_recharged").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("企业累计充值金额"),

		// 关联字段
		field.Int64("admin_user_id").
			Comment("企业管理员对应的 users.id"),
	}
}

func (Enterprise) Edges() []ent.Edge {
	return []ent.Edge{
		// 企业管理员 —— 指向 users 表
		edge.From("admin", User.Type).
			Ref("managed_enterprises").
			Field("admin_user_id").
			Required().
			Unique(),

		// 企业成员
		edge.To("members", EnterpriseMember.Type).
			Comment("企业下的所有成员记录"),

		// 企业部门
		edge.To("departments", Department.Type).
			Comment("企业下的部门"),

		// 企业套餐
		edge.To("subscriptions", EnterpriseSubscription.Type).
			Comment("企业购买的套餐记录"),

		// 企业用量日志
		edge.To("usage_logs", UsageLog.Type).
			Comment("企业维度的用量日志"),
	}
}

func (Enterprise) Indexes() []ent.Index {
	return []ent.Index{
		// name 唯一约束（WHERE deleted_at IS NULL）由迁移 SQL 实现部分索引
		index.Fields("name"),
		index.Fields("parent_id"),
		index.Fields("status"),
		index.Fields("admin_user_id"),
		index.Fields("deleted_at"),
	}
}
