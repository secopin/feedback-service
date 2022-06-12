package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Feedback struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (feedback *Feedback) Format() string {
	value := reflect.ValueOf(feedback).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() == "" {
			value.Field(i).SetString("-")
		}
	}

	return fmt.Sprintf(
		"<b>Отправитель:</b> %s\n<b>Компания:</b> %s\n<b>Телефон:</b> %s\n<b>Почта:</b> %s\n\n<b>Сообщение:</b> %s",
		feedback.Name, feedback.Company, feedback.Phone, feedback.Email, feedback.Message,
	)
}

func main() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("FB_TOKEN"))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var feedback *Feedback

		err := json.NewDecoder(request.Body).Decode(&feedback)
		if err != nil {
			writer.WriteHeader(400)
			_, _ = writer.Write([]byte("{\"status\":\"error while processing a request\"}"))
			return
		}

		payload := strings.NewReader(fmt.Sprintf(
			"{\"text\":\"%s\",\"parse_mode\":\"HTML\",\"chat_id\":\"%s\"}",
			feedback.Format(), os.Getenv("FB_CHANNEL")),
		)

		telegram, err := http.Post(url, "application/json", payload)
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("{\"status\":\"error while sending internal request\"}"))
			return
		}
		if telegram.StatusCode != 200 {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("{\"status\":\"internal query returned unexpected result\"}"))
			return
		}

		response, err := json.Marshal(feedback)
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("{\"status\":\"error while generating response\"}"))
			return
		}

		writer.WriteHeader(200)
		_, _ = writer.Write(response)
		return
	})

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("FB_PORT")), nil)
	if err != nil {
		panic(err)
	}
}
