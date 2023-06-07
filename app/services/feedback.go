package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type FeedbackService struct{}

func (f FeedbackService) PassFeedback(feedback structs.Feedback) error {
	// if !f.checkGoogleCaptcha(feedback.Token) {
	// 	return structs.ErrCheckCaptcha
	// }

	fmt.Println(f.checkGoogleCaptcha(feedback.Token))

	return nil
}

func (f FeedbackService) checkGoogleCaptcha(token string) bool {
	var googleCaptcha string = "6Lc0EvwlAAAAAFosqn77174Ehz966vtEKz-XszRp"

	req, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)

	pathQuery := req.URL.Query()
	pathQuery.Add("secret", googleCaptcha)
	pathQuery.Add("response", token)

	req.URL.RawQuery = pathQuery.Encode()

	client := &http.Client{}

	var googleResponse map[string]interface{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &googleResponse)
	return googleResponse["success"].(bool)
}

func (f FeedbackService) send(name, email, comment string) {
}
