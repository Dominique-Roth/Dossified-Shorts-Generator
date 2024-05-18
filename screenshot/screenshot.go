package screenshot

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

    "go.uber.org/zap"

	"github.com/Dominique-Roth/Dossified-Shorts-Generator/config"
	"github.com/Dominique-Roth/Dossified-Shorts-Generator/logging"
	"github.com/Dominique-Roth/Dossified-Shorts-Generator/rest"
	"github.com/Dominique-Roth/Dossified-Shorts-Generator/utils"
)

func ScreenshotTrends(trends []rest.TrendArticle) {
    logging.Info("Taking screenshots")
	configuration := config.GetConfiguration()
	gowitnessHost := configuration.GowitnessHost
	restApiHost := configuration.RemoteUrl

	screenshotPostUrl := gowitnessHost + "/api/screenshot"
	screenshotPath := createScreenshotDir()

	for _, article := range trends {
		requestUrl := fmt.Sprintf(
			"%s/vid_gen/?item_id=%d&obj_type=%s",
			restApiHost,
			article.ArticleId,
			article.ArticleType,
		)
        logging.Debug("url", zap.String("url", requestUrl))
		body := []byte(`{
            "url": "` + requestUrl + `",
            "oneshot": "true"
        }`)

		// Create a HTTP post request
		screenshotRequest, err := http.NewRequest("POST", screenshotPostUrl, bytes.NewBuffer(body))
		utils.CheckError(err)
		screenshotRequest.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(screenshotRequest)
		utils.CheckError(err)
		defer res.Body.Close()

		// Save screenshot to png file
		file, err := os.Create(screenshotPath + "/" + fmt.Sprint(article.ArticleId) + ".png")
		utils.CheckError(err)
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		utils.CheckError(err)
	}
    logging.Info("Screenshots taken")
}

func createScreenshotDir() string {
	path := filepath.Join(".", "output/screenshots")
	err := os.MkdirAll(path, os.ModePerm)
	utils.CheckError(err)
	logging.Debug("Screenshot path", zap.String("path", path))
	return path
}
