package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	mixin "graphQlDemo/mixins"
)

// Review holds the schema definition for the Review entity.
type Review struct {
	ent.Schema
}

// Fields of the Review.
func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Int("rating"),
		field.String("comment"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Review.
func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("reviewer", User.Type).
			Ref("reviews").
			Unique(),
		edge.From("reviwedTool", Tool.Type).
			Ref("reviews").
			Unique(),
	}
}

func (Review) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate()),
    }
}

func (Review) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ToolTime{},
	}
}
