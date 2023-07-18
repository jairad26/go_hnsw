package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func weaviate_run() ([][]float64, []string) {
	cfg := weaviate.Config{
		Host:       "example", // Replace with your endpoint
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: "example"}, // Replace w/ your Weaviate instance API key
		Headers: map[string]string{
			"X-HuggingFace-Api-Key": "example", // Replace with your inference API key
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	className := "Question"

	// delete the existing question class
	if err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()); err != nil {
		// Weaviate will return a 400 if the class does not exist, so this is allowed, only return an error if it's not a 400
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode != http.StatusBadRequest {
			panic(err)
		}
	}

	classObj := &models.Class{
		Class:      "Question",
		Vectorizer: "text2vec-huggingface", // If set to "none" you must always provide vectors yourself. Could be any other "text2vec-*" also.
		ModuleConfig: map[string]interface{}{
			"text2vec-huggingface": map[string]interface{}{
				"model": "sentence-transformers/paraphrase-MiniLM-L6-v2",
				"options": map[string]interface{}{
					"waitForModel": true,
				},
			},
		},
	}

	// add the schema
	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// Retrieve the data
	items, err := getJSONdata()
	if err != nil {
		panic(err)
	}

	// convert items into a slice of models.Object
	objects := make([]*models.Object, len(items))
	for i := range items {
		objects[i] = &models.Object{
			Class: "Question",
			Properties: map[string]any{
				"category": items[i]["Category"],
				"question": items[i]["Question"],
				"answer":   items[i]["Answer"],
			},
		}
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			panic(res.Result.Errors.Error)
		}
	}

	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
		{Name: "_additional { id vector }"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		// WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	var vec_arr [][]float64
	var question_arr []string
	for i := 0; i < 10; i++ {
		question := result.Data["Get"].(map[string]interface{})[className].([]interface{})[i].(map[string]interface{})["question"].(string)
		vec_interface := result.Data["Get"].(map[string]interface{})[className].([]interface{})[i].(map[string]interface{})["_additional"].(map[string]interface{})["vector"].([]interface{})
		vec := make([]float64, len(vec_interface))
		for i := 0; i < len(vec_interface); i++ {
			vec[i] = vec_interface[i].(float64)
		}
		// fmt.Printf("%T", vec)
		// pp.Print(vec)
		// fmt.Println("-----------------------------------------")
		vec_arr = append(vec_arr, vec)
		question_arr = append(question_arr, question)
	}

	// formatted_result, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pp.Print(result)
	return vec_arr, question_arr
}

func getJSONdata() ([]map[string]string, error) {
	// Retrieve the data
	data, err := http.DefaultClient.Get("https://raw.githubusercontent.com/weaviate-tutorials/quickstart/main/data/jeopardy_tiny.json")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()

	// Decode the data
	var items []map[string]string
	if err := json.NewDecoder(data.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
