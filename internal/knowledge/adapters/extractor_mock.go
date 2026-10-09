package adapters

func entityChatResponseMock() entityChatResponse {
	return entityChatResponse{
		Choices: []struct {
			Message entityChatMessage `json:"message"`
		}{
			{
				Message: entityChatMessage{
					Role:    "assistant",
					Content: "{\"entities\": [[\"Minestrone Recipe\"], [\"Ingredients\", \"olive oil\", \"onion\", \"celery\", \"carrot\", \"courgette\", \"smoked pancetta\", \"garlic\", \"dried oregano\", \"cannellini beans\", \"chopped tomatoes\", \"tomato purée\", \"vegetable stock\", \"bay leaf\", \"pasta\", \"greens\", \"kale\", \"chard\", \"cavolo nero\", \"basil\", \"parmesan\"], [\"Method\"], [\"onion\", \"celery\", \"carrot\", \"courgette\", \"pancetta\", \"garlic\", \"oregano\", \"beans\", \"chopped tomatoes\", \"purée\", \"stock\", \"bay leaf\"], [\"pasta\", \"greens\", \"basil\", \"parmesan\"], [\"BBC Good Food\", \"Classic Minestrone Soup\"]]}",
				},
			},
		},
		Error: nil,
	}
}
