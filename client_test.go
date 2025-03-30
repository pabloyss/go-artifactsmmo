package artifactsmmo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
	"github.com/stretchr/testify/assert"
)

// mockHTTPClient creates a mock HTTP client that returns predefined responses
func mockHTTPClient(t *testing.T, expectedMethod, expectedPath string, statusCode int, response interface{}) *http.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedPath, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if response != nil {
			responseData, err := json.Marshal(response)
			assert.NoError(t, err)
			_, err = w.Write(responseData)
			assert.NoError(t, err)
		}
	}))

	t.Cleanup(func() {
		server.Close()
	})

	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		},
	}
}

func TestNewClient(t *testing.T) {
	client := artifactsmmo.NewClient("test-token", "testuser")
	assert.NotNil(t, client)
	assert.Equal(t, "test-token", client.Config.GetToken())
	assert.Equal(t, "testuser", client.Config.GetUsername())
}

func TestNewClientWithCustomHttpClient(t *testing.T) {
	httpClient := &http.Client{}
	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)
	assert.NotNil(t, client)
	assert.Equal(t, httpClient, client.Config.GetClient())
}

func TestFight(t *testing.T) {
	expectedResponse := &models.CharacterFight{}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/fight",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Fight()

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetCharacterInfo(t *testing.T) {
	expectedResponse := &models.Character{
		Name:      "testuser",
		Level:     10,
		Health:    100,
		MaxHealth: 100,
		Gold:      500,
		Exp:       1234,
		MaxExp:    2000,
	}

	httpClient := mockHTTPClient(t, "GET", "/characters/testuser",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.GetCharacterInfo("testuser")

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestMove(t *testing.T) {
	expectedResponse := &models.CharacterMovementData{
		From: models.Coords{X: 10, Y: 10},
		To:   models.Coords{X: 11, Y: 11},
		Map:  "forest",
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/move",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Move(11, 11)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestMove_AlreadyAtDestination(t *testing.T) {
	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/move",
		490, nil)

	client := NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Move(10, 10)

	assert.Error(t, err)
	assert.Equal(t, models.ErrAlreadyAtDestination, err)
	assert.Nil(t, result)
}

func TestEquip(t *testing.T) {
	expectedResponse := &models.EquipRequest{
		Item: models.Item{
			Code: "sword",
			Name: "Iron Sword",
			Type: models.ItemTypeWeapon,
		},
		Slot: models.SlotWeapon,
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/equip",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Equip("sword", models.SlotWeapon, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestUnequip(t *testing.T) {
	expectedResponse := &models.EquipRequest{
		Item: models.Item{
			Code: "sword",
			Name: "Iron Sword",
			Type: models.ItemTypeWeapon,
		},
		Slot: models.SlotWeapon,
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/unequip",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Unequip(models.SlotWeapon, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGather(t *testing.T) {
	expectedResponse := &models.SkillData{
		Exp:     50,
		Level:   5,
		Name:    "Mining",
		Type:    models.SkillTypeMining,
		Success: true,
		Items:   []models.Drop{{Code: "iron_ore", Quantity: 1}},
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/gathering",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Gather()

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetItem(t *testing.T) {
	expectedResponse := &models.SingleItem{
		Item: models.Item{
			Code:     "health_potion",
			Name:     "Health Potion",
			Type:     models.ItemTypeConsumable,
			Level:    1,
			BuyPrice: 50,
		},
	}

	httpClient := mockHTTPClient(t, "GET", "/items/health_potion",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.GetItem("health_potion")

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetItem_NotFound(t *testing.T) {
	httpClient := mockHTTPClient(t, "GET", "/items/nonexistent",
		http.StatusNotFound, nil)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.GetItem("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, models.ErrItemNotFound, err)
	assert.Nil(t, result)
}

func TestAcceptNewTask(t *testing.T) {
	expectedResponse := &models.TaskData{
		Task: models.Task{
			Code:     "mining_task",
			Name:     "Mine 10 Iron Ore",
			Progress: 0,
			Goal:     10,
			Type:     models.TaskTypeGathering,
		},
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/task/new",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.AcceptNewTask()

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestCompleteTask(t *testing.T) {
	expectedResponse := &models.TaskRewardData{
		Exp:   100,
		Gold:  50,
		Coins: 2,
		Items: []models.Drop{{Code: "iron_ore", Quantity: 5}},
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/task/complete",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.CompleteTask()

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestRest(t *testing.T) {
	expectedResponse := &models.Rest{
		Health: 100,
		Gold:   -10,
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/rest",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.Rest()

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestUseItem(t *testing.T) {
	expectedResponse := &models.UseItem{
		Gold:   0,
		Health: 25,
		Effect: "Restored 25 health points",
	}

	httpClient := mockHTTPClient(t, "POST", "/my/testuser/action/use",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.UseItem("health_potion", 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetMap(t *testing.T) {
	expectedResponse := &models.MapSchema{
		X:    10,
		Y:    10,
		Name: "Forest",
		Type: "forest",
	}

	httpClient := mockHTTPClient(t, "GET", "/maps/10/10",
		http.StatusOK, expectedResponse)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.GetMap(10, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetMap_NotFound(t *testing.T) {
	httpClient := mockHTTPClient(t, "GET", "/maps/999/999",
		http.StatusNotFound, nil)

	client := artifactsmmo.NewClientWithCustomHttpClient("test-token", "testuser", httpClient)

	result, err := client.GetMap(999, 999)

	assert.Error(t, err)
	assert.Equal(t, models.ErrMapNotFound, err)
	assert.Nil(t, result)
}
