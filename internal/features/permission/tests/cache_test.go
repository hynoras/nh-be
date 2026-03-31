package permission

import (
	"context"
	"errors"
	"nh-be/internal/features/permission"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCache_GetCodeNames(t *testing.T) {
	userId := TestUserID()
	key := "perm:" + userId.String()

	tests := []struct {
		name           string
		userId         uuid.UUID
		setupMock      func(mock redismock.ClientMock)
		expectedResult []string
		expectedError  error
	}{
		{
			name:   "success_with_code_names",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(key).SetVal(TestCodeNames())
			},
			expectedResult: TestCodeNames(),
			expectedError:  nil,
		},
		{
			name:   "success_with_no_members",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(key).SetVal([]string{})
			},
			expectedResult: []string{},
			expectedError:  nil,
		},
		{
			name:   "success_when_key_does_not_exist",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(key).RedisNil()
			},
			expectedResult: nil,
			expectedError:  redis.Nil,
		},
		{
			name:   "redis_returns_error",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(key).SetErr(redis.ErrClosed)
			},
			expectedResult: nil,
			expectedError:  redis.ErrClosed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setupMock(mock)
			cache := permission.NewPermissionCache(db)

			result, err := cache.GetCodeNames(context.Background(), tc.userId)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCache_SetCodeNames(t *testing.T) {
	userId := TestUserID()
	key := "perm:" + userId.String()

	tests := []struct {
		name          string
		userId        uuid.UUID
		codeNames     []string
		setupMock     func(mock redismock.ClientMock)
		expectedError error
	}{
		{
			name:      "empty_code_names",
			userId:    userId,
			codeNames: []string{},
			setupMock: func(mock redismock.ClientMock) {
				// No mock expectations — Redis should not be called
			},
			expectedError: nil,
		},
		{
			name:      "delete_existing_key",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
		{
			name:      "add_all_members_to_set",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
		{
			name:      "set_expiration",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
		{
			name:      "pipeline_exec_success",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
		{
			name:      "pipeline_exec_failure",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetErr(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
		{
			name:      "handle_duplicate_values",
			userId:    userId,
			codeNames: TestDuplicateCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "READ", "READ").SetVal(1)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
		{
			name:      "correct_key_generation",
			userId:    userId,
			codeNames: TestCodeNames(),
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
				mock.ExpectSAdd(key, "view_experiment", "manage_experiment").SetVal(2)
				mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setupMock(mock)
			cache := permission.NewPermissionCache(db)

			err := cache.SetCodeNames(context.Background(), tc.userId, tc.codeNames)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCache_InvalidateUser(t *testing.T) {
	userId := TestUserID()
	key := "perm:" + userId.String()

	tests := []struct {
		name          string
		userId        uuid.UUID
		setupMock     func(mock redismock.ClientMock)
		expectedError error
	}{
		{
			name:   "success",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(1)
			},
			expectedError: nil,
		},
		{
			name:   "success_key_not_exist",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetVal(0)
			},
			expectedError: nil,
		},
		{
			name:   "redis_error",
			userId: userId,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectDel(key).SetErr(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setupMock(mock)
			cache := permission.NewPermissionCache(db)

			err := cache.InvalidateUser(context.Background(), tc.userId)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCache_InvalidateAll(t *testing.T) {
	errScan := errors.New("scan error")
	errDel := errors.New("del error")
	errDel1 := errors.New("del error 1")
	errDel2 := errors.New("del error 2")

	scanKey := "perm:*"

	tests := []struct {
		name       string
		setupMock  func(mock redismock.ClientMock)
		checkError func(t *testing.T, err error)
	}{
		{
			name: "success_single_scan_iteration",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1", "perm:u2"}, 0)
				mock.ExpectDel("perm:u1", "perm:u2").SetVal(2)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success_multiple_scan_iterations",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 10)
				mock.ExpectDel("perm:u1").SetVal(1)
				mock.ExpectScan(10, scanKey, 100).SetVal([]string{"perm:u2"}, 0)
				mock.ExpectDel("perm:u2").SetVal(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success_no_keys_found",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{}, 0)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success_with_empty_scan_batches",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{}, 10)
				mock.ExpectScan(10, scanKey, 100).SetVal([]string{"perm:u1"}, 0)
				mock.ExpectDel("perm:u1").SetVal(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "success_cursor_ends_immediately",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 0)
				mock.ExpectDel("perm:u1").SetVal(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "fail_when_single_del_fails",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 0)
				mock.ExpectDel("perm:u1").SetErr(errDel)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errDel)
			},
		},
		{
			name: "fail_when_multiple_del_failures",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 10)
				mock.ExpectDel("perm:u1").SetErr(errDel1)
				mock.ExpectScan(10, scanKey, 100).SetVal([]string{"perm:u2"}, 0)
				mock.ExpectDel("perm:u2").SetErr(errDel2)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errDel1)
				assert.ErrorIs(t, err, errDel2)
			},
		},
		{
			name: "success_when_del_fails_but_other_batches_succeed",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 10)
				mock.ExpectDel("perm:u1").SetVal(1)
				mock.ExpectScan(10, scanKey, 100).SetVal([]string{"perm:u2"}, 0)
				mock.ExpectDel("perm:u2").SetErr(errDel)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errDel)
			},
		},
		{
			name: "fail_when_scan_fails_immediately",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetErr(errScan)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errScan)
			},
		},
		{
			name: "fail_when_scan_fails_mid_iteration",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 10)
				mock.ExpectDel("perm:u1").SetVal(1)
				mock.ExpectScan(10, scanKey, 100).SetErr(errScan)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errScan)
			},
		},
		{
			name: "fail_with_del_and_scan_errors_combined",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1"}, 10)
				mock.ExpectDel("perm:u1").SetErr(errDel)
				mock.ExpectScan(10, scanKey, 100).SetErr(errScan)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errDel)
				assert.ErrorIs(t, err, errScan)
			},
		},
		{
			name: "success_no_errors_returns_nil",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, scanKey, 100).SetVal([]string{"perm:u1", "perm:u2"}, 0)
				mock.ExpectDel("perm:u1", "perm:u2").SetVal(2)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setupMock(mock)
			cache := permission.NewPermissionCache(db)

			err := cache.InvalidateAll(context.Background())

			tc.checkError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
