package testing

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"circle-2.0/app/pkg/repository/user"
	"circle-2.0/lib/connection/database"
	"circle-2.0/lib/helper/constant"
	"circle-2.0/lib/model"
	"circle-2.0/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	app := test.SetUpTestApp()

	requestBody := model.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	requestBodyBytes, err := app.Config().JSONEncoder(requestBody)
	if err != nil {
		assert.Nil(t, err)
	}

	req := httptest.NewRequest(http.MethodPost, constant.CreateUserPath, bytes.NewReader(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Nil(t, err)
	}

	var responseData model.Response
	if err = app.Config().JSONDecoder(respBody, &responseData); err != nil {
		assert.Nil(t, err)
	}

	actualUserData, ok := responseData.Data.(model.Users)
	if !ok {
		assert.Equal(t, true, ok)
	}

	u := user.NewUserRepo(database.MySQLConnect)
	
	userData, err := u.GetUserByEmail(requestBody.Email)
	if err != nil {
		assert.Nil(t, err)
    }

	if err = u.DeleteUser(userData.ID); err != nil {
		assert.Nil(t, err)
    }

	assert.Equal(t, requestBody.Name, actualUserData.Name)
	assert.Equal(t, requestBody.Email, actualUserData.Email)
	assert.Equal(t, requestBody.Password, actualUserData.Password)
	assert.NotEmpty(t, actualUserData.ID)
	assert.NotEmpty(t, actualUserData.CreatedAt)
	assert.NotEmpty(t, actualUserData.UpdatedAt)
	assert.Empty(t, actualUserData.UpdatedBy)
}
