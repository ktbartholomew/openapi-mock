package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ktbartholomew/openapi-mock/config"
	"github.com/ktbartholomew/openapi-mock/template"

	"gopkg.in/yaml.v3"
)

// OpenAPISpec is the entire document for an API
type OpenAPISpec struct {
	Paths map[string]PathSpec
}

// PathSpec is the definition of a single URI path
type PathSpec struct {
	Delete *MethodSpec `yaml:"delete,omitempty"`
	Get    *MethodSpec `yaml:"get,omitempty"`
	Patch  *MethodSpec `yaml:"patch,omitempty"`
	Post   *MethodSpec `yaml:"post,omitempty"`
}

// MethodSpec is the definition of a method for a given path
type MethodSpec struct {
	Summary   string
	Responses map[string]ResponseSpec `yaml:"responses"`
}

// ResponseSpec is the definition of one of an endpoint's possible responses
type ResponseSpec struct {
	Description string
	Content     ContentSpec `json:"content"`
}

// ContentSpec is the definition of the contents an endpoint might return
type ContentSpec struct {
	JSONContent *ContentTypeSpec `yaml:"application/json"`
}

// ContentTypeSpec is the definition of a single content payload that an endpoint might return
type ContentTypeSpec struct {
	Example string `yaml:"example"`
}

func main() {
	config.Setup()
	cfg := config.Get()
	file, err := getSpecFile(cfg)
	if err != nil {
		panic(err)
	}

	r, err := createRouter(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("listening on %s...\n", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createRouter(spec []byte) (*mux.Router, error) {
	r := mux.NewRouter()

	doc := OpenAPISpec{}
	yaml.Unmarshal(spec, &doc)

	for path, spec := range doc.Paths {
		fmt.Printf("adding path %s\n", path)

		if spec.Get != nil {
			r.Methods(http.MethodGet).Path(path).HandlerFunc(responderFunc(spec.Get.Responses))
		}

		if spec.Patch != nil {
			r.Methods(http.MethodPatch).Path(path).HandlerFunc(responderFunc(spec.Patch.Responses))
		}

		if spec.Post != nil {
			r.Methods(http.MethodPost).Path(path).HandlerFunc(responderFunc(spec.Post.Responses))
		}

		if spec.Delete != nil {
			r.Methods(http.MethodDelete).Path(path).HandlerFunc(responderFunc(spec.Delete.Responses))
		}
	}

	return r, nil
}

func responderFunc(responses map[string]ResponseSpec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status, resp, err := getRequestedOrDefaultResponse(responses, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		mockLatency(r)

		w.Header().Set("X-Mock-Description", resp.Description)
		if resp.Content.JSONContent != nil {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)

		if resp.Content.JSONContent == nil {
			return
		}

		count, err := strconv.Atoi(r.Header.Get("X-Mock-Count"))
		if err != nil {
			count = 1
		}

		td := template.TemplateData{
			Params:    mux.Vars(r),
			ItemCount: count,
		}

		body, err := template.ExampleOutput(resp.Content.JSONContent.Example, td)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(body)
		return

	}
}

func getSpecFile(cfg *config.Config) ([]byte, error) {
	if cfg.SpecPath != "" {
		file, err := ioutil.ReadFile(cfg.SpecPath)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	if cfg.SpecURL != "" {
		resp, err := http.Get(cfg.SpecURL)
		if err != nil {
			return nil, err
		}

		bb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bb, nil
	}

	return nil, fmt.Errorf("neither a path or URL to an OpenAPI spec file was provided")
}

func getRequestedOrDefaultResponse(responses map[string]ResponseSpec, request *http.Request) (int, ResponseSpec, error) {
	// If a specific response has been requested, try to return that
	if request.Header.Get("X-Mock-Response") != "" {
		requestedResponse := request.Header.Get("X-Mock-Response")

		if _, ok := responses[requestedResponse]; !ok {
			return 0, ResponseSpec{}, fmt.Errorf("no response available for requested status code")
		}

		i, err := strconv.Atoi(requestedResponse)
		if err != nil {
			return 0, ResponseSpec{}, err
		}

		return i, responses[requestedResponse], nil
	}

	// Otherwise, fall back to the response with the lowest status code
	statuses := []int{}
	for k := range responses {
		if k == "default" {
			continue
		}

		i, _ := strconv.Atoi(k)
		statuses = append(statuses, i)
	}

	return statuses[0], responses[strconv.Itoa(statuses[0])], nil
}

func mockLatency(r *http.Request) {
	if latstring := r.Header.Get("X-Mock-Latency"); latstring != "" {
		latency, err := strconv.Atoi(latstring)
		if err != nil {
			time.Sleep(0)
		} else {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}
}
