package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ToolTime holds the schema definition for the ToolTime entity.
type ToolTime struct {
	mixin.Schema
}

// Fields of the ToolTime.
func (ToolTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Immutable().
			Default(func() time.Time {
				return time.Now().UTC()
			}).Annotations(entgql.OrderField("CREATE_TIME")).Comment("The time this object was created at"),
		field.Time("update_time").
			Default(func() time.Time {
				return time.Now().UTC()
			}).
			UpdateDefault(func() time.Time {
				return time.Now().UTC()
			}).Annotations(entgql.OrderField("UPDATE_TIME")).Comment("The last time this object was mutated"),
	}
}
