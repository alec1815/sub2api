package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/internal/domain"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Department holds the schema definition for the Department entity.
// 企业组织架构——部门表。
type Department struct {
	ent.Schema
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "departments"},
	}
}

func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("enterprise_id").
			Comment("所属企业 ID"),

		// 注意：parent_id 直接作为字段使用，不定义 edge
		// 因为 parent_id=0 表示顶级部门，无对应记录
		field.Int64("parent_id").
			Default(0).
			Comment("父部门 ID，0=顶级部门"),

		field.String("name").
			MaxLen(100).
			NotEmpty().
			Comment("部门名称"),

		field.Int("order_num").
			Default(0).
			Comment("排序号，越小越靠前"),

		// 负责人信息
		field.String("leader").
			MaxLen(100).
			Default("").
			Comment("负责人姓名"),
		field.String("phone").
			MaxLen(50).
			Default("").
			Comment("负责人电话"),
		field.String("email").
			MaxLen(255).
			Default("").
			Comment("负责人邮箱"),

		// 状态
		field.String("status").
			MaxLen(20).
			Default(domain.StatusActive).
			Comment("active / disabled"),
	}
}

func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		// 所属企业
		edge.From("enterprise", Enterprise.Type).
			Ref("departments").
			Field("enterprise_id").
			Required().
			Unique(),

		// 部门下的成员
		edge.To("members", EnterpriseMember.Type).
			Comment("该部门下的成员"),
	}
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		// (enterprise_id, name) 唯一（WHERE deleted_at IS NULL）—— 同一企业下名称唯一
		// 部分唯一索引由迁移 SQL 实现
		index.Fields("enterprise_id", "name"),
		// 查子部门
		index.Fields("enterprise_id", "parent_id"),
		index.Fields("deleted_at"),
	}
}
