package ExpretionData

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
)

var (
	ExprationDataNotFound = errors.New("Expration data is not found.")
	PronunciationNotFound = errors.New("Pronunciation data is not found.")
	NoTranslationFound    = errors.New("Translation data is not found.")
)

const (
	ExampleApiKey     = "8b2ebc45-6b7c-4575-89be-537ca676f94a"
	TranslationApiKey = "53e30d15-e741-4a64-bb35-43ae0c36745e:fx"
	RequestAttemts    = 3
)

type Pronunciation struct {
	Phonetic string
	Path     string
}

type Translation struct {
	Translation string
	Examples    []string
}

type ExpretionData struct {
	Translations  []Translation
	Pronunciation Pronunciation
}

type Response struct {
	Hwi struct {
		Prs []struct {
			Mw    string `json:"mw"`
			Sound struct {
				Audio string `json:"audio"`
			} `json:"sound"`
		} `json:"prs"`
	} `json:"hwi"`
	Def []struct {
		Sseq [][][]interface{} `json:"sseq"`
	} `json:"def"`
}

type DeepLResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func GetEpretionData(expretion string, requestAttemts int) (ExpretionData, error) {
	url := fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/learners/json/%s?key=%s", expretion, ExampleApiKey)

	resp, err := http.Get(url)
	if err != nil {
		if requestAttemts != 0 {
			return GetEpretionData(expretion, requestAttemts-1)
		} else {
			return ExpretionData{}, ExprationDataNotFound
		}
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var data []Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return ExpretionData{}, ExprationDataNotFound
	}

	// Check if data is empty
	if len(data) == 0 {
		return ExpretionData{}, ExprationDataNotFound
	}

	var examples []string

	for _, def := range data[0].Def {
		for _, sseq := range def.Sseq {
			for _, sense := range sseq {
				if len(sense) > 1 {
					senseData, ok := sense[1].(map[string]interface{})
					if !ok {
						continue
					}
					if dt, ok := senseData["dt"].([]interface{}); ok {
						for _, item := range dt {
							itemData, ok := item.([]interface{})
							if !ok || len(itemData) < 2 {
								continue
							}
							if itemData[0] == "vis" {
								for _, vis := range itemData[1].([]interface{}) {
									visData, ok := vis.(map[string]interface{})
									if !ok {
										continue
									}
									if example, ok := visData["t"].(string); ok {
										examples = append(examples, example)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	var max int
	if len(examples) < 8 {
		max = len(examples)
	} else {
		max = 8
	}

	examples = correctExamples(examples[:max])

	var translations []Translation

	for _, context := range examples {
		contextTranslation, err := getTransltion(expretion, context, requestAttemts)

		if err != nil {
			continue
		}
		existingContextTranslationIndex := slices.IndexFunc(translations, func(translation Translation) bool {
			return translation.Translation == contextTranslation
		})

		if existingContextTranslationIndex == -1 {
			translations = append(translations, Translation{contextTranslation, []string{context}})
		} else {
			translations[existingContextTranslationIndex].Examples = append(translations[existingContextTranslationIndex].Examples, context)
		}
	}

	return ExpretionData{translations, getPronuciation(expretion, data)}, nil
}

func getPronuciation(expretion string, data []Response) Pronunciation {
	var pronunciation Pronunciation
	for _, pr := range data[0].Hwi.Prs {
		if pr.Mw != "" {
			pronunciation.Phonetic = pr.Mw
		}
		if pr.Sound.Audio != "" {
			audioURL := fmt.Sprintf("https://media.merriam-webster.com/soundc11/%s/%s.wav", string(pr.Sound.Audio[0]), pr.Sound.Audio)
			downloadFile(expretion, audioURL)
			pronunciation.Path = fmt.Sprintf("/home/nazar/nazzar/vocbl/audio/%v", expretion)
		}
	}

	return pronunciation

}

func downloadFile(expretion string, url string) error {
	// Create the file
	out, err := os.Create(fmt.Sprintf("/home/nazar/nazzar/vocbl/audio/%v", expretion))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func correctExamples(examples []string) []string {
	var corrextedExamples = make([]string, len(examples))
	re := regexp.MustCompile(`\{[^}]*\}`)
	for i, example := range examples {
		corrextedExamples[i] = re.ReplaceAllString(example, "")
	}
	return corrextedExamples
}

func getTransltion(expretion, context string, tries int) (string, error) {

	apiURL := "https://api-free.deepl.com/v2/translate"

	data := url.Values{}
	data.Set("auth_key", TranslationApiKey)
	data.Set("text", expretion)
	data.Set("source_lang", "EN")
	data.Set("target_lang", "UK")
	data.Set("context", context)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		if tries != 0 {
			return getTransltion(expretion, context, tries-1)
		} else {

			return "", NoTranslationFound
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var deepLResponse DeepLResponse
	if err := json.Unmarshal(bodyBytes, &deepLResponse); err != nil {
		fmt.Println("Error decoding the response:", err)
		os.Exit(1)
	}

	if len(deepLResponse.Translations) > 0 {

		return deepLResponse.Translations[0].Text, nil
	} else {

		return deepLResponse.Translations[0].Text, NoTranslationFound
	}

}
