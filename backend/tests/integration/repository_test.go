package integration

import (
	"testing"

	"github.com/labhaus/backend/internal/domain/style"
	"github.com/labhaus/backend/internal/infrastructure/persistence"
)

func TestStyleRepository_Integration(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	repo := persistence.NewStyleRepository(testDB.DB)
	ctx := ClearContext()

	t.Run("Create and FindByID", func(t *testing.T) {
		// Create style - New(name, description, prompt, category, tags)
		entity, err := style.New("Anime Style", "Japanese anime art", "anime, colorful", "Art", []string{"anime", "japanese"})
		if err != nil {
			t.Fatalf("Failed to create entity: %v", err)
		}

		err = repo.Create(ctx, entity)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if entity.ID == "" {
			t.Error("ID should be generated")
		}

		// Find by ID
		found, err := repo.FindByID(ctx, entity.ID)
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}

		if found.Name != "Anime Style" {
			t.Errorf("Name = %v, want Anime Style", found.Name)
		}
		if len(found.Tags) != 2 {
			t.Errorf("Tags length = %d, want 2", len(found.Tags))
		}
	})

	t.Run("Update", func(t *testing.T) {
		// Create style
		entity, _ := style.New("Original", "desc", "prompt", "Cat", []string{"tag1"})
		repo.Create(ctx, entity)

		// Update
		entity.Name = "Updated"
		entity.Description = "updated desc"
		err := repo.Update(ctx, entity)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify
		found, _ := repo.FindByID(ctx, entity.ID)
		if found.Name != "Updated" {
			t.Errorf("Name = %v, want Updated", found.Name)
		}
		if found.Description != "updated desc" {
			t.Errorf("Description = %v, want updated desc", found.Description)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		// Create multiple styles
		s1, _ := style.New("Style1", "desc1", "prompt1", "Art", []string{"tag1"})
		s2, _ := style.New("Style2", "desc2", "prompt2", "Photo", []string{"tag2"})
		s3, _ := style.New("Style3", "desc3", "prompt3", "Art", []string{"tag3"})
		repo.Create(ctx, s1)
		repo.Create(ctx, s2)
		repo.Create(ctx, s3)

		// FindAll
		all, err := repo.FindAll(ctx, style.Filter{})
		if err != nil {
			t.Fatalf("FindAll() error = %v", err)
		}

		if len(all) < 3 {
			t.Errorf("FindAll() got %d styles, want at least 3", len(all))
		}
	})

	t.Run("FindAll with category filter", func(t *testing.T) {
		// Create style in specific category
		entity, _ := style.New("Cyberpunk Style", "Futuristic neon art", "cyberpunk, neon, futuristic", "SciFi", []string{"cyberpunk"})
		repo.Create(ctx, entity)

		// Find by category using Filter
		found, err := repo.FindAll(ctx, style.Filter{Category: "SciFi"})
		if err != nil {
			t.Fatalf("FindAll(Category) error = %v", err)
		}

		if len(found) == 0 {
			t.Error("FindAll(Category) should return at least 1 style")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Create style
		entity, _ := style.New("ToDelete", "desc", "prompt", "Cat", []string{})
		repo.Create(ctx, entity)

		// Delete
		err := repo.Delete(ctx, entity.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify deleted
		_, err = repo.FindByID(ctx, entity.ID)
		if err == nil {
			t.Error("FindByID() should return error for deleted style")
		}
	})
}

// FIXME: This test is misplaced - it's trying to test User with style.New
// Should be removed or moved to user repository tests
func TestUserRepository_Integration(t *testing.T) {
	t.Skip("Skipping misplaced test - style.New cannot create users")
	
	// testDB := SetupTestDB(t)
	// defer testDB.Close()
	// defer testDB.Cleanup(t)
	
	// repo := persistence.NewStyleRepository(testDB.DB)
	// ctx := ClearContext()
	
	// t.Run("should work", func(t *testing.T) {
	// 	// This is wrong - style.New() cannot create a user
	// 	// entity, err := style.New("test@example.com", "hashedpassword123", "Test User", "user")
	// })
}
