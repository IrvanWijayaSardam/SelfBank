package mocks

import (
	"testing"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testData := []entity.Deposit{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
	}

	mockRepo.On("All", 1, 10).Return(testData, nil)

	result, err := mockRepo.All(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, testData, result)

	mockRepo.AssertExpectations(t)
}

func TestFindDepositByID(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testID := uint64(1)
	testData := &entity.Deposit{ID: testID, Amount: 100}

	mockRepo.On("FindDepositByID", testID).Return(testData)

	result := mockRepo.FindDepositByID(testID)

	assert.NotNil(t, result)
	assert.Equal(t, testData, result)

	mockRepo.AssertExpectations(t)
}

func TestFindDepositByIDUser(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testID := uint64(1)
	page := 1
	pageSize := 10
	testData := []entity.Deposit{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
	}

	mockRepo.On("FindDepositByIDUser", testID, page, pageSize).Return(testData, nil)
	result, err := mockRepo.FindDepositByIDUser(testID, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, testData, result)

	mockRepo.AssertExpectations(t)
}

func TestInsertDeposit(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testDeposit := &entity.Deposit{ID: 1, Amount: 100}

	mockRepo.On("InsertDeposit", testDeposit).Return(*testDeposit)

	result := mockRepo.InsertDeposit(testDeposit)

	assert.NotNil(t, result)
	assert.Equal(t, *testDeposit, result)

	mockRepo.AssertExpectations(t)
}

func TestStorePaymentToken(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testTransactionID := uint64(1)
	testPaymentToken := "payment_token"
	testVirtualAcc := "virtual_acc"

	mockRepo.On("StorePaymentToken", testTransactionID, testPaymentToken, testVirtualAcc).Return(nil)

	err := mockRepo.StorePaymentToken(testTransactionID, testPaymentToken, testVirtualAcc)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTotalDeposit(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testTotal := int64(1000)

	mockRepo.On("TotalDeposit").Return(testTotal)

	result := mockRepo.TotalDeposit()

	assert.Equal(t, testTotal, result)

	mockRepo.AssertExpectations(t)
}

func TestTotalDepositByUserID(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testUserID := uint64(1)
	testTotal := int64(500)

	mockRepo.On("TotalDepositByUserID", testUserID).Return(testTotal)

	result := mockRepo.TotalDepositByUserID(testUserID)

	assert.Equal(t, testTotal, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDeposit(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testDeposit := entity.Deposit{ID: 1, Amount: 200}

	mockRepo.On("UpdateDeposit", testDeposit).Return(testDeposit)

	result := mockRepo.UpdateDeposit(testDeposit)

	assert.Equal(t, testDeposit, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDepositStatus(t *testing.T) {
	mockRepo := NewDepositRepository(t)

	testID := uint64(1)
	testNewStatus := uint64(2)

	mockRepo.On("UpdateDepositStatus", testID, testNewStatus).Return(nil)

	err := mockRepo.UpdateDepositStatus(testID, testNewStatus)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
