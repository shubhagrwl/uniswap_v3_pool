package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	poolDTO "uniswapper/internal/app/db/dto/pool"
	testutils "uniswapper/internal/app/service/util/testutils/mocks"
	mockDB "uniswapper/internal/app/service/util/testutils/mocks/repository/pool"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) {
	envPath := "../../../../.env"
	testutils.SetupTest(t, envPath)
}

func TestGetPoolLogsById(t *testing.T) {
	setupTest(t)

	poolID := "123"
	blockID := "34522"

	testCases := []struct {
		name          string
		url           string
		buildStubs    func(store *mockDB.MockIPoolLogsRepository)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "status ok 202",
			url:  fmt.Sprintf("/pool/%s?block=%s", poolID, blockID),
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {

				store.
					EXPECT().
					GetPoolLogs(gomock.Any(), poolID, blockID).
					Return(poolDTO.Logs{PoolAddress: poolID}, nil).
					Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, resp.Code)
			},
		},
		{
			name: "status 500",
			url:  fmt.Sprintf("/pool/%s?block=%s", poolID, blockID),
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {

				store.
					EXPECT().
					GetPoolLogs(gomock.Any(), poolID, blockID).
					Return(poolDTO.Logs{PoolAddress: poolID}, fmt.Errorf("error while getting info")).
					Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "bad request 400",
			url:  "/pool/ ",
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPoolDBClient := mockDB.NewMockIPoolLogsRepository(ctrl)
			tc.buildStubs(mockPoolDBClient)

			controller := NewPoolController(mockPoolDBClient)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/pool/:pool_id", controller.GetPoolLogsById)

			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tc.checkResponse(t, resp)
		})
	}
}

func TestGetHistoryPoolLogs(t *testing.T) {
	setupTest(t)

	poolID := "123"

	testCases := []struct {
		name          string
		url           string
		buildStubs    func(store *mockDB.MockIPoolLogsRepository)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "status ok 202",
			url:  fmt.Sprintf("/pool/%s/historic", poolID),
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {

				store.
					EXPECT().
					GetPoolLogsHistory(gomock.Any(), poolID).
					Return([]poolDTO.Logs{{PoolAddress: poolID}}, nil).
					Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, resp.Code)
			},
		},
		{
			name: "status 500",
			url:  fmt.Sprintf("/pool/%s/historic", poolID),
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {

				store.
					EXPECT().
					GetPoolLogsHistory(gomock.Any(), poolID).
					Return([]poolDTO.Logs{{PoolAddress: poolID}}, fmt.Errorf("error while getting info")).
					Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "bad request 400",
			url:  "/pool//historic",
			buildStubs: func(store *mockDB.MockIPoolLogsRepository) {
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPoolDBClient := mockDB.NewMockIPoolLogsRepository(ctrl)
			tc.buildStubs(mockPoolDBClient)

			controller := NewPoolController(mockPoolDBClient)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/pool/:pool_id/historic", controller.GetPoolLogsHistory)

			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tc.checkResponse(t, resp)
		})
	}
}
