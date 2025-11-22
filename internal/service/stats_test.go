package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStatsRepo struct {
	mock.Mock
}

func (m *MockStatsRepo) GetAssignedReviewersCountPerPR(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	v, ok := args.Get(0).(map[string]int)
	if !ok {
		return nil, args.Error(1)
	}
	return v, args.Error(1)
}

func (m *MockStatsRepo) GetOpenPRCountPerUser(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	v, ok := args.Get(0).(map[string]int)
	if !ok {
		return nil, args.Error(1)
	}
	return v, args.Error(1)
}

func TestStatsService_Getters(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockStatsRepo)

	svc := NewStatsService(mockRepo)

	expectedPRs := map[string]int{"pr1": 2, "pr2": 0}
	expectedUsers := map[string]int{"u1": 1, "u2": 3}

	mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(expectedPRs, nil)
	mockRepo.On("GetOpenPRCountPerUser", mock.Anything).Return(expectedUsers, nil)

	prRes, err := svc.GetAssignedCountPerPR(ctx)
	assert.NoError(t, err)
	if !reflect.DeepEqual(prRes, expectedPRs) {
		t.Fatalf("unexpected pr result: got=%v want=%v", prRes, expectedPRs)
	}

	userRes, err := svc.GetOpenPRCountPerUser(ctx)
	assert.NoError(t, err)
	if !reflect.DeepEqual(userRes, expectedUsers) {
		t.Fatalf("unexpected user result: got=%v want=%v", userRes, expectedUsers)
	}

	mockRepo.AssertExpectations(t)
}

func TestStatsService_ErrorPropagation(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockStatsRepo)
	svc := NewStatsService(mockRepo)

	mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(nil, errors.New("db error"))
	mockRepo.On("GetOpenPRCountPerUser", mock.Anything).Return(nil, errors.New("db error"))

	_, err := svc.GetAssignedCountPerPR(ctx)
	assert.Error(t, err)

	_, err = svc.GetOpenPRCountPerUser(ctx)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
