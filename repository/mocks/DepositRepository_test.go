// repository_test.go

package mocks

import (
	"testing"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testData := []entity.Deposit{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
	}

	// Set expectations for the mock method
	mockRepo.On("All", 1, 10).Return(testData, nil)

	// Call the method you're testing (replace with your actual repository method)
	result, err := mockRepo.All(1, 10)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, testData, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestFindDepositByID(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testID := uint64(1)
	testData := &entity.Deposit{ID: testID, Amount: 100}

	// Set expectations for the mock method
	mockRepo.On("FindDepositByID", testID).Return(testData)

	// Call the method you're testing (replace with your actual repository method)
	result := mockRepo.FindDepositByID(testID)

	// Assert the results
	assert.NotNil(t, result)
	assert.Equal(t, testData, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestFindDepositByIDUser(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testID := uint64(1)
	page := 1
	pageSize := 10
	testData := []entity.Deposit{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
	}

	// Set expectations for the mock method
	mockRepo.On("FindDepositByIDUser", testID, page, pageSize).Return(testData, nil)

	// Call the method you're testing (replace with your actual repository method)
	result, err := mockRepo.FindDepositByIDUser(testID, page, pageSize)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, testData, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestInsertDeposit(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testDeposit := &entity.Deposit{ID: 1, Amount: 100}

	// Set expectations for the mock method
	mockRepo.On("InsertDeposit", testDeposit).Return(*testDeposit)

	// Call the method you're testing (replace with your actual repository method)
	result := mockRepo.InsertDeposit(testDeposit)

	// Assert the results
	assert.NotNil(t, result)
	assert.Equal(t, *testDeposit, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestStorePaymentToken(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testTransactionID := uint64(1)
	testPaymentToken := "payment_token"
	testVirtualAcc := "virtual_acc"

	// Set expectations for the mock method
	mockRepo.On("StorePaymentToken", testTransactionID, testPaymentToken, testVirtualAcc).Return(nil)

	// Call the method you're testing (replace with your actual repository method)
	err := mockRepo.StorePaymentToken(testTransactionID, testPaymentToken, testVirtualAcc)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestTotalDeposit(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testTotal := int64(1000)

	// Set expectations for the mock method
	mockRepo.On("TotalDeposit").Return(testTotal)

	// Call the method you're testing (replace with your actual repository method)
	result := mockRepo.TotalDeposit()

	// Assert the results
	assert.Equal(t, testTotal, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestTotalDepositByUserID(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testUserID := uint64(1)
	testTotal := int64(500)

	// Set expectations for the mock method
	mockRepo.On("TotalDepositByUserID", testUserID).Return(testTotal)

	// Call the method you're testing (replace with your actual repository method)
	result := mockRepo.TotalDepositByUserID(testUserID)

	// Assert the results
	assert.Equal(t, testTotal, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeposit(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testDeposit := entity.Deposit{ID: 1, Amount: 200}

	// Set expectations for the mock method
	mockRepo.On("UpdateDeposit", testDeposit).Return(testDeposit)

	// Call the method you're testing (replace with your actual repository method)
	result := mockRepo.UpdateDeposit(testDeposit)

	// Assert the results
	assert.Equal(t, testDeposit, result)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdateDepositStatus(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := NewDepositRepository(t)

	// Define the test data and expected result
	testID := uint64(1)
	testNewStatus := uint64(2)

	// Set expectations for the mock method
	mockRepo.On("UpdateDepositStatus", testID, testNewStatus).Return(nil)

	// Call the method you're testing (replace with your actual repository method)
	err := mockRepo.UpdateDepositStatus(testID, testNewStatus)

	// Assert the results
	assert.NoError(t, err)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}
