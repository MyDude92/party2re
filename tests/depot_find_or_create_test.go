package party2re_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corecharacter "github.com/witchcraze/party2re/internal/character"
	"github.com/witchcraze/party2re/internal/depot"
)

// MockDepotRepo verifies depot.FindOrCreate calls across domain callers
type MockDepotRepo struct {
	mock.Mock
}

func (m *MockDepotRepo) FindByCharacterIDForUpdate(ctx context.Context, characterID string) (depot.Depot, error) {
	args := m.Called(ctx, characterID)
	return args.Get(0).(depot.Depot), args.Error(1)
}

func (m *MockDepotRepo) Save(ctx context.Context, d depot.Depot) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func TestDepotFindOrCreate_ExistingDepot(t *testing.T) {
	mockRepo := new(MockDepotRepo)
	ctx := context.Background()
	char := corecharacter.Character{
		ID:        "char_test_101",
		JobLevel:  50,
		OverDepot: 2,
	}

	existingDepot := depot.Depot{
		CharacterID: char.ID,
		Capacity:    10,
		ExDepot:     0,
	}

	mockRepo.On("FindByCharacterIDForUpdate", ctx, char.ID).Return(existingDepot, nil)

	d, err := depot.FindOrCreate(ctx, mockRepo, char)
	assert.NoError(t, err)
	assert.Equal(t, char.ID, d.CharacterID)
	// RefreshCapacity should apply JobLevel and OverDepot
	assert.GreaterOrEqual(t, d.Capacity, 10)
	mockRepo.AssertExpectations(t)
}

func TestDepotFindOrCreate_UninitializedDepotFallback(t *testing.T) {
	mockRepo := new(MockDepotRepo)
	ctx := context.Background()
	char := corecharacter.Character{
		ID:        "char_uninit_202",
		JobLevel:  10,
		OverDepot: 0,
	}

	// Returns ErrNotFound for uninitialized character
	mockRepo.On("FindByCharacterIDForUpdate", ctx, char.ID).Return(depot.Depot{}, depot.ErrNotFound)

	d, err := depot.FindOrCreate(ctx, mockRepo, char)
	assert.NoError(t, err)
	assert.Equal(t, char.ID, d.CharacterID)
	assert.GreaterOrEqual(t, d.Capacity, 5)
	mockRepo.AssertExpectations(t)
}
