package wecom_robot

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

)

type WeComRobot struct {
	Key string
}

func NewWeComRobot (key string) *WeComRobot{
	return &WeComRobot{
		Key: key,
	}
}
func (w *WeComRobot) Notice(ctx context.Context, content string) error {
	res, err := http.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", w.Key),
		"application/json",
		strings.NewReader(
			fmt.Sprintf(`
	{
		"msgtype": "text",
		"text": {
		"content": "%s"
	}
	}
	`, content)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (w *WeComRobot) SendFile(ctx context.Context, filepath string, filename string) error {
	mediaID, err := w.uploadToWecom(filepath, filename)
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", w.Key),
		"application/json",
		strings.NewReader(
			fmt.Sprintf(`
	{
		"msgtype": "file",
		"file": {
		"media_id": "%s"
	}
	}
	`, mediaID)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func (w *WeComRobot) uploadToWecom(filepath string, filename string) (mediaID string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	res, err := w.UploadFile(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file", w.Key),
		nil, "media", filename,
		f)
	if err != nil {
		return "", err
	}
	var ret map[string]interface{}
	_ = JSONUnmarshal(res, &ret)
	if ret["media_id"] == nil {
		return "", errors.New("media_id 不存在")
	}
	return ret["media_id"].(string), nil
}
func (w *WeComRobot) UploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	HTTPClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func JSONUnmarshal(data []byte, v interface{}) error {
	buffer := bytes.NewBuffer(data)
	decoder := json.NewDecoder(buffer)
	decoder.UseNumber()
	return decoder.Decode(&v)
}

func ToCsvRow(rows ...string) string {
	for i, v := range rows {
		rows[i] = ReplaceComma(v)
	}
	return strings.Join(rows, ",") + "\n"
}

func MustAppendFile(filePath string, data string) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	write := bufio.NewWriter(f)
	_, _ = write.WriteString(data)
	_ = write.Flush()
}

func ReplaceComma(s string) string {
	return strings.ReplaceAll(s, ",", "，")
}