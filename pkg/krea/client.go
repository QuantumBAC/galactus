package krea

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GenImageOpts struct {
	ImageResolution             string  `json:"image-resolution"`
	PromptInfluence             int     `json:"prompt-influence"`
	DiffusionSteps              int     `json:"diffusion-steps"`
	ControlnetConditioningScale float64 `json:"controlnet-conditioning-scale"`
	PatternURL                  float64 `json:"pattern-url"`
	ImageURL                    string  `json:"image-url"`
	Prompt                      string  `json:"prompt"`
	NegativePrompt              string  `json:"negative-prompt"`
}

type genImageResponse struct {
	Data string `json:"data"`
}

const baseURL = "https://canvas.krea.ai"

type Client interface {
	GenImage(opts GenImageOpts) (string, error)
}

type kreaImpl struct {
	httpClient *http.Client
}

func New() Client {
	return &kreaImpl{
		httpClient: &http.Client{},
	}
}

// GenImage passes GenImageOpts to the PatternizeFun endpoint, and returns the url stored in the `data` field of the response object if the call succeeds
func (k *kreaImpl) GenImage(opts GenImageOpts) (string, error) {

	// Create the request URL with appropriate query parameters.
	endpoint := "/?%2FpatternizeFun"
	u, _ := url.Parse(baseURL + endpoint)

	// Create a buffer to write the form data.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fmt.Printf("Generating image with options: %#v\n", opts)

	// Add form fields to the request.
	writer.WriteField("image-resolution", opts.ImageResolution)
	writer.WriteField("prompt-influence", strconv.Itoa(opts.PromptInfluence))
	writer.WriteField("diffusion-steps", strconv.Itoa(opts.DiffusionSteps))
	writer.WriteField("controlnet-conditioning-scale", fmt.Sprintf("%f", opts.ControlnetConditioningScale))
	writer.WriteField("pattern-url", fmt.Sprintf("%f", opts.PatternURL))
	writer.WriteField("image-url", opts.ImageURL)
	writer.WriteField("prompt", opts.Prompt)

	if opts.NegativePrompt != "" {
		writer.WriteField("negative-prompt", opts.NegativePrompt)
	}

	writer.Close()

	req, err := http.NewRequest("POST", u.String(), &requestBody)
	if err != nil {
		return "", err
	}

	// Set the Content-Type header to the appropriate multipart content type.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-SvelteKit-Action", "true")
	req.Header.Set("Origin", baseURL)

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var response genImageResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return "", err
	}

	return findImageURLInData([]byte(response.Data))

	// fmt.Printf("%#v %s %s \n", response, err, url)

	// return url, err
}

func findImageURLInData(data []byte) (string, error) {
	var l []interface{}

	if err := json.Unmarshal(data, &l); err != nil {
		return "", err
	}

	for _, e := range l {
		if url, ok := e.(string); ok && strings.HasPrefix(url, "https://canvas-generations") {
			return url, nil
		}
	}

	return "", errors.New(fmt.Sprintf("no url found in response body: %#v", l))
}
