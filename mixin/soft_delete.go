package mixin

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"
	"fmt"
)

var _ ent.Mixin = (*SoftDelete)(nil)

type SoftDelete struct{ mixin.Schema }

func (SoftDelete) Fields() []ent.Field {
	return DeletedAt{}.Fields()
}

func (SoftDelete) Hooks() []ent.Hook {
	return []ent.Hook{
		PrintHook,
	}
}

func PrintHook(next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		fmt.Printf("Type: %s, Operation: %s, ConcreteType: %T\n", m.Type(), m.Op(), m)
		return next.Mutate(ctx, m)
	})
}
