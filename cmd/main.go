package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"test_entgo_soft_delete/ent"
	_ "test_entgo_soft_delete/ent/runtime"
	"test_entgo_soft_delete/ent/user"
)

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.Debug().User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func UpdateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.Debug().User.
		UpdateOneID(1).
		SetAge(18).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed update user: %w", err)
	}
	log.Println("user was updated: ", u)
	return u, nil
}

func DeleteUser(ctx context.Context, client *ent.Client) error {
	return client.Debug().User.
		DeleteOneID(1).
		Exec(ctx)
}

func SoftDeleteUser(ctx context.Context, client *ent.Client) error {
	_, err := client.Debug().User.
		UpdateOneID(1).
		SetDeletedAt(time.Now()).
		Save(ctx)
	return err
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.Debug().User.
		Query().
		Where(user.NameEQ("a8m")).
		Where(user.AgeEQ(18)).
		Where(user.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateEntClient(ctx context.Context) *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

func main() {
	ctx := context.Background()
	var err error

	client := CreateEntClient(ctx)
	defer client.Close()

	if _, err = CreateUser(ctx, client); err != nil {
		log.Fatal(err)
	}
	if _, err = QueryUser(ctx, client); err != nil {
		log.Fatal(err)
	}
	if err = DeleteUser(ctx, client); err != nil {
		log.Fatal(err)
	}
}
