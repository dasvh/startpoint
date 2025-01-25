package loader

import (
	goccy "github.com/goccy/go-yaml"
	"gopkg.in/yaml.v3"
	sigyaml "sigs.k8s.io/yaml"
	"testing"
)

var largeYamlData = []byte(`
prev_req:
url: foobar.com
method: POST
headers:
  X-Foo-Bar: SomeValue
  Content-Type: application/json
auth:
  basic:
    username: user
    password: pw
body: >
  {
    "id": 1,
    "name": "Jane",
    "details": {
      "age": 30,
      "address": {
        "street": "123 Main St",
        "city": "Anytown",
        "zipcode": "12345"
      },
      "phones": ["123-456-7890", "987-654-3210"]
    }
  }
`)

type Request struct {
	PrevReq string `yaml:"prev_req"`
	Url     string `yaml:"url"`
	Method  string `yaml:"method"`
	Headers struct {
		ContentType string `yaml:"Content-Type"`
	} `yaml:"headers"`
	Auth struct {
		Basic struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basic"`
	} `yaml:"auth"`
	Body string `yaml:"body"`
}

func BenchmarkSigK8sUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var req Request
		if err := sigyaml.Unmarshal(largeYamlData, &req); err != nil {
			b.Fatalf("failed to unmarshal with sig.k8s.io/yaml: %v", err)
		}
	}
}

func BenchmarkGoccyUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var req Request
		if err := goccy.Unmarshal(largeYamlData, &req); err != nil {
			b.Fatalf("failed to unmarshal with goccy/go-yaml: %v", err)
		}
	}
}

func BenchmarkYamlV3Unmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var req Request
		if err := yaml.Unmarshal(largeYamlData, &req); err != nil {
			b.Fatalf("failed to unmarshal with gopkg.in/yaml.v3: %v", err)
		}
	}
}

func BenchmarkSigK8sMarshal(b *testing.B) {
	req := Request{
		PrevReq: "prev_req",
		Url:     "foobar.com",
		Method:  "POST",
		Headers: struct {
			ContentType string `yaml:"Content-Type"`
		}{ContentType: "application/json"},
		Auth: struct {
			Basic struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"basic"`
		}{Basic: struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{Username: "user", Password: "pw"}},
		Body: `{
			"id": 1,
			"name": "Jane",
			"details": {
				"age": 30,
				"address": {
					"street": "123 Main St",
					"city": "Anytown",
					"zipcode": "12345"
				},
				"phones": ["123-456-7890", "987-654-3210"]
			}
		}`,
	}
	for i := 0; i < b.N; i++ {
		if _, err := sigyaml.Marshal(&req); err != nil {
			b.Fatalf("failed to marshal with sig.k8s.io/yaml: %v", err)
		}
	}
}

func BenchmarkGoccyMarshal(b *testing.B) {
	req := Request{
		PrevReq: "prev_req",
		Url:     "foobar.com",
		Method:  "POST",
		Headers: struct {
			ContentType string `yaml:"Content-Type"`
		}{ContentType: "application/json"},
		Auth: struct {
			Basic struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"basic"`
		}{Basic: struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{Username: "user", Password: "pw"}},
		Body: `{
			"id": 1,
			"name": "Jane",
			"details": {
				"age": 30,
				"address": {
					"street": "123 Main St",
					"city": "Anytown",
					"zipcode": "12345"
				},
				"phones": ["123-456-7890", "987-654-3210"]
			}
		}`,
	}
	for i := 0; i < b.N; i++ {
		if _, err := goccy.Marshal(&req); err != nil {
			b.Fatalf("failed to marshal with goccy/go-yaml: %v", err)
		}
	}
}

func BenchmarkYamlV3Marshal(b *testing.B) {
	req := Request{
		PrevReq: "prev_req",
		Url:     "foobar.com",
		Method:  "POST",
		Headers: struct {
			ContentType string `yaml:"Content-Type"`
		}{ContentType: "application/json"},
		Auth: struct {
			Basic struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"basic"`
		}{Basic: struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{Username: "user", Password: "pw"}},
		Body: `{
			"id": 1,
			"name": "Jane",
			"details": {
				"age": 30,
				"address": {
					"street": "123 Main St",
					"city": "Anytown",
					"zipcode": "12345"
				},
				"phones": ["123-456-7890", "987-654-3210"]
			}
		}`,
	}
	for i := 0; i < b.N; i++ {
		if _, err := yaml.Marshal(&req); err != nil {
			b.Fatalf("failed to marshal with gopkg.in/yaml.v3: %v", err)
		}
	}
}
