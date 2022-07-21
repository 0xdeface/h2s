package main

import (
	"emailer/internal/message_maker"
	"encoding/json"
	"fmt"
	"html/template"
)

func main() {
	m := message_maker.MessageMaker{}
	var data interface{}
	jsonRaw := `
		{
        "name": "Banana",
        "points": 200,
        "description": "A banana grown in Ecuador."
    }`
	err := json.Unmarshal([]byte(jsonRaw), &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, _ := template.New("t1").Parse("sample template {{ .name }}")
	result, _ := m.Make(t, data)
	fmt.Println(string(result))
}
