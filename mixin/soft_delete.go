package mixin

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"time"
)

var _ ent.Mixin = (*SoftDelete)(nil)

type SoftDelete struct{ mixin.Schema }

func (SoftDelete) Fields() []ent.Field {
	return DeletedAt{}.Fields()
}

func (SoftDelete) Hooks() []ent.Hook {
	return []ent.Hook{
		SoftDeleteHook,
	}
}

func PrintHook(next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		fmt.Printf("Type: %s, Operation: %s, ConcreteType: %T\n", m.Type(), m.Op(), m)
		return next.Mutate(ctx, m)
	})
}

func SoftDeleteHook(next ent.Mutator) ent.Mutator {
	type DeletedAtSetter interface {
		SetDeletedAt(time.Time)
	}
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		ns, ok := m.(DeletedAtSetter)
		if !ok {
			return nil, fmt.Errorf("unexpected soft-delete call from mutation type %T", m)
		}

		switch op := m.Op(); {
		case op.Is(ent.OpDelete | ent.OpDeleteOne):
			ns.SetDeletedAt(time.Now())
		}

		return next.Mutate(ctx, m)
	})
}
