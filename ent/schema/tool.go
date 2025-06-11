package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	mixin "graphQlDemo/mixins"
)

// Tool holds the schema definition for the Tool entity.
type Tool struct {
	ent.Schema
}

// Fields of the Tool.
func (Tool) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").Unique(),
		field.String("description"),
		field.Enum("category").Values(
			"FRONTEND",
			"BACKEND",
			"FULLSTACK",
			"MOBILE",
			"DEVOPS",
			"TESTING",
			"DATABASE",
			"CLOUD",
			"SECURITY",
			"MONITORING",
			"VERSION_CONTROL",
			"DOCUMENTATION",
		),
		field.String("website"),
		field.String("image_url"),
		field.Float("average_rating").Default(0).Optional(),
        field.Int("rating_count").Default(0).Optional(),
	}
}

// Edges of the Tool.
func (Tool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("reviews", Review.Type),
	}
}

func (Tool) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate()),
    }
}

func (Tool) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ToolTime{},
	}
}
