package config

type Endpoint struct {
	Name    string
	Address string
}

type Config struct {
	Endpoints []Endpoint
}

func Default() Config {
	return Config{
		Endpoints: []Endpoint{
			{
				Name:    "Yandex",
				Address: "https://ya.ru",
			},
			{
				Name:    "MAX",
				Address: "https://max.ru",
			},
			{
				Name:    "VK",
				Address: "https://vk.com",
			},
			{
				Name:    "RuStore",
				Address: "https://rustore.ru",
			},
		},
	}
}
