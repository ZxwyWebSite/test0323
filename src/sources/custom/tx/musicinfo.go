package tx

import (
	"encoding/gob"
	"lx-source/src/env"
	"lx-source/src/sources"
	"strings"

	"github.com/ZxwyWebSite/ztool"
	"github.com/ZxwyWebSite/ztool/x/bytesconv"
)

func init() {
	gob.Register(musicInfo{})
}

func getMusicInfo(songMid string) (infoBody musicInfo, emsg string) {
	cquery := strings.Join([]string{`tx`, songMid, `info`}, `/`)
	if cdata, ok := env.Cache.Get(cquery); ok {
		if cinfo, ok := cdata.(musicInfo); ok {
			infoBody = cinfo
			return
		}
	}
	infoReqBody := ztool.Str_FastConcat(`{"comm":{"ct":"19","cv":"1859","uin":"0"},"req":{"method":"get_song_detail_yqq","module":"music.pf_song_detail_svr","param":{"song_mid":"`, songMid, `","song_type":0}}}`)
	var infoResp struct {
		Code int `json:"code"`
		// Ts      int64  `json:"ts"`
		// StartTs int64  `json:"start_ts"`
		// Traceid string `json:"traceid"`
		Req struct {
			Code int       `json:"code"`
			Data musicInfo `json:"data"`
		} `json:"req"`
	}
	err := signRequest(bytesconv.StringToBytes(infoReqBody), &infoResp)
	if err != nil {
		emsg = err.Error()
		return //nil, err.Error()
	}
	if infoResp.Code != 0 || infoResp.Req.Code != 0 {
		emsg = `获取音乐信息失败`
		return //nil, `获取音乐信息失败`
	}
	infoBody = infoResp.Req.Data
	env.Cache.Set(cquery, infoBody, sources.C_lx)
	return //infoBody.Req.Data, ``
}
