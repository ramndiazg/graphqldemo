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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name"),
		field.String("username").Unique(),
		field.String("email").Unique(),
		field.String("password_hash"),
		field.Enum("role").Values("user", "admin").Default("user"),
		field.Bool("is_verified").Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("reviews", Review.Type),
	}
}

func (User) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate()),
    }
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.ToolTime{},
	}
}
