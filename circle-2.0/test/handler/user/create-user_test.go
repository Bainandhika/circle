package testing

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"circle-2.0/app/config"
	"circle-2.0/app/pkg/repository/user"
	"circle-2.0/lib/connection/database"
	"circle-2.0/lib/connection/nosql"
	"circle-2.0/lib/helper/constant"
	"circle-2.0/lib/model"
	"circle-2.0/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	app := test.SetUpTestApp()
	defer func ()  {
		db, _ := database.MySQLConnect.DB()
		db.Close()

		nosql.RedisConnect.Close()
	}()

	requestBody := model.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	requestBodyBytes, err := app.Config().JSONEncoder(requestBody)
	if err != nil {
		assert.Nil(t, err, err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodPost, constant.CreateUserPath, bytes.NewReader(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("secret-key", config.App.SecretKey)
	req.Header.Set("channel-id", "test-channel")

	resp, err := app.Test(req, -1)
	if err != nil {
		assert.Nil(t, err)
		return
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Nil(t, err, err.Error())
		return
	}

	var responseData model.Response
	if err = app.Config().JSONDecoder(respBody, &responseData); err != nil {
		assert.Nil(t, err, err.Error())
		return
	}

	u := user.NewUserRepo(database.MySQLConnect)

	userData, err := u.GetUserByEmail(requestBody.Email)
	if err != nil {
		assert.Nil(t, err, err.Error())
		return
	}

	if err = u.DeleteUser(userData.ID); err != nil {
		assert.Nil(t, err, err.Error())
		return
	}

	assert.Equal(t, requestBody.Name, userData.Name)
	assert.Equal(t, requestBody.Email, userData.Email)
	assert.NotEmpty(t, userData.Password)
	assert.NotEmpty(t, userData.ID)
	assert.NotEmpty(t, userData.CreatedAt)
	assert.NotEmpty(t, userData.UpdatedAt)
	assert.NotEmpty(t, userData.UpdatedBy)
}
