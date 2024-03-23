// 网络请求相关 (beta); 分类: Net_(网络)

package ztool

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	stdurl "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZxwyWebSite/ztool/x/json"
	"github.com/ZxwyWebSite/ztool/x/maps"
)

type (
	Net_ReqHandlerFunc func(req *http.Request) error
	Net_ResHandlerFunc func(res *http.Response) error

	// Net_HandlerErr struct {
	// 	p string
	// 	e error
	// }
)

const (
	// 默认User-Agent
	Net_ua = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36 HBPC/12.1.2.300`
	// 默认Accept
	Net_accept = `application/json, text/plain, */*`
)

var (
	// 默认Client
	Net_client = &http.Client{
		Timeout: 30 * time.Second,
		// Transport: &http.Transport{
		// 	DisableKeepAlives: true,
		// },
	}
	// 默认Header
	Net_header = map[string]string{
		`User-Agent`: Net_ua,
		`Accept`:     Net_accept,
	}
	// 定义一些报错
	Net_ErrNotA301    = errors.New(` Not A 301 Resp`)
	Net_ErrDownOut    = errors.New(`no output filename and cant infer from url`)
	Net_ErrNoRedirect = errors.New(`no 301 or 302 status code returned`)
	// go:lint disable-last-line
)

// HttpRequest 网络请求 (方法, 地址, 数据, Header, JSON) (原数据, 错误)
func Net_HttpReq(method, url string, body io.Reader, header map[string]string, out any) ([]byte, error) {
	// client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set(`User-Agent`, Net_ua)
	req.Header.Add(`Accept`, Net_accept)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := Net_client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if out != nil {
		json.Unmarshal(data, out)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

// func HandlerAddHeader(h map[string]string) func(*http.Request) {
// 	return func(r *http.Request) {
// 		for k, v := range h {
// 			r.Header.Add(k, v)
// 		}
// 	}
// }

// Get301url 获取重定向后的地址 (地址) (重定向地址, 错误)
func Net_HttpGet301(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(`User-Agent`, Net_ua)
	resp, err := client.Do(req)
	if err != nil {
		return ``, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 301 && resp.StatusCode != 302 {
		return ``, Net_ErrNotA301
	}
	return resp.Header.Get(`Location`), nil
}

// ==================== ↓通用框架重构版↓ ====================

// 通用框架
func request(client *http.Client, method, url string, body io.Reader, reqh []Net_ReqHandlerFunc, resh []Net_ResHandlerFunc) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	for i, r := 0, len(reqh); i < r; i++ {
		if err = reqh[i](req); err != nil {
			return err
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	for i, r := 0, len(resh); i < r; i++ {
		if err = resh[i](res); err != nil {
			return err
		}
	}
	return nil
}

// 批量添加请求头
func Net_ReqAddHeader(headers map[string]string) Net_ReqHandlerFunc {
	return func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Add(k, v)
			// req.Header.Add(k, headers[k])
		}
		return nil
	}
}
func Net_ReqAddHeaders(headers ...map[string]string) Net_ReqHandlerFunc {
	buf_header := maps.Clone(Net_header)
	for i := range headers { // v: map[string]string
		if headers[i] != nil {
			maps.Copy(buf_header, headers[i])
		}
	}
	return Net_ReqAddHeader(buf_header)
}

// 映射到结构体
// func Net_ResToStruct(out any) Net_ResHandlerFunc {
// 	return func(res *http.Response) (err error) {
// 		var data []byte
// 		data, err = io.ReadAll(res.Body)
// 		if err == nil {
// 			err = json.Unmarshal(data, out)
// 			if err != nil {
// 				err = fmt.Errorf(`%s, data: %s`, err, data)
// 			}
// 		}
// 		return
// 	}
// }

// JSON body 映射到结构体 (流式解析测试)
func Net_ResToStruct(out any) Net_ResHandlerFunc {
	return func(res *http.Response) error {
		dec := json.NewDecoder(res.Body)
		// dec.UseNumber()
		return dec.Decode(out)
	}
}

// 流式下载
func Net_Download(url string, out io.Writer, header map[string]string) error {
	return request(Net_client, http.MethodGet, url, nil,
		[]Net_ReqHandlerFunc{
			Net_ReqAddHeaders(header),
		},
		[]Net_ResHandlerFunc{
			func(res *http.Response) error {
				// buf := make([]byte, 4*1024) // 4kb (对应http默认缓存) //128kb
				// _, err := io.CopyBuffer(out, res.Body, buf)
				// buf = nil // 手动释放缓存，防止内存泄露
				_, err := io.Copy(out, res.Body)
				return err
			},
		})
}
func Net_DownloadFile(url, out string, header map[string]string) error {
	if out == `` {
		ubj, err := stdurl.Parse(url)
		if err != nil {
			return err
		}
		if name := filepath.Base(ubj.Path); name != `.` && name != string(os.PathSeparator) {
			out = name
		} else {
			return Net_ErrDownOut // 未提供输出文件名且无法从Url推断
		}
	}
	file, err := Fbj_CreatFile(out)
	// if err != nil {
	// 	return err
	// }
	if err == nil {
		err = Net_Download(url, file, header)
		file.Close()
	}
	return err
}

// 获取重定向后的地址
func Net_GetRedirectAddr(url string, header map[string]string) (out string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	err = request(client, http.MethodGet, url, nil,
		[]Net_ReqHandlerFunc{
			Net_ReqAddHeaders(header),
		},
		[]Net_ResHandlerFunc{
			func(res *http.Response) error {
				// StatusCode = 301 or 302
				if res.StatusCode == http.StatusMovedPermanently || res.StatusCode == http.StatusFound {
					out = res.Header.Get(`Location`)
					return nil
				}
				return Net_ErrNoRedirect // 没有返回重定向
			},
		})
	return
}

// 发起多功能请求
func Net_Request(method, url string, body io.Reader, reqh []Net_ReqHandlerFunc, resh []Net_ResHandlerFunc) error {
	return request(Net_client, method, url, body, reqh, resh)
}

// 创建表单
func Net_Values(val map[string]string) string {
	buf := stdurl.Values{}
	for k, v := range val {
		buf.Set(k, v)
	}
	return buf.Encode()
}
func Net_FormData(val map[string]string) io.Reader {
	return strings.NewReader(Net_Values(val))
}

// WebKitFormBoundary <form> <body, content-type>
func Net_MultiPart(val map[string]string) (io.Reader, string) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	for k, v := range val {
		mw.WriteField(k, v)
	}
	mw.Close()
	return buf, mw.FormDataContentType()
}
