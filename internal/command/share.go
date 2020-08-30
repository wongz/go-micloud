package command

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
	"go-micloud/pkg/zlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (r *Command) Share() *cli.Command {
	return &cli.Command{
		Name:  "share",
		Usage: "Get public share url",
		Action: func(context *cli.Context) error {
			var args = context.Args()
			for i := 0; i < args.Len(); i++ {
				fileName := args.Get(i)
				fileName = strings.ReplaceAll(fileName, "\\s", " ")
				fileInfo, ok := FileMap[fileName]
				if !ok {
					return errors.New("当前目录不存在该文件")
				}
				if fileInfo.Type == "folder" {
					return errors.New("目前不支持分享文件夹")
				}
				downloadUrl, err := r.HttpApi.GetFileDownLoadUrl(fileInfo.Id)
				if err != nil {
					return errors.New(fmt.Sprintf("获取下载地址失败：%s", err.Error()))
				}
				var shortUrl = downloadUrl
				resp, err := http.PostForm("http://t.wibliss.com/api/v1/create", url.Values{"url": []string{downloadUrl}})
				if err == nil {
					all, _ := ioutil.ReadAll(resp.Body)
					dataUrl := gjson.Get(string(all), "data.url").String()
					if dataUrl != "" {
						shortUrl = dataUrl
					}
					resp.Body.Close()
				}
				zlog.Logger.Info(shortUrl)

				i := strings.LastIndex(shortUrl, "/")

				shortUrl = shortUrl[:i] + "?t=" + shortUrl[i+1:]

				zlog.Info(fmt.Sprintf("获取分享成功,有效期24小时，复制链接( %s )到浏览器里面打开下载,请注意浏览器弹框！", shortUrl))
			}
			return nil
		},
	}
}