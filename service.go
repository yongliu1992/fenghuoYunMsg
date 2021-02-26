package fenghuoYunMsg

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type SmsSendService struct {
	Req PostReq
}

type PostReq struct {
	Account     string `json:"account"`
	Token       string `json:"token"`
	Ts          string `json:"ts"`
	SignatureId string `json:"signatureId"`
	Mobiles     string `json:"mobiles"`
	Content     string `json:"content"`
	Ref         string `json:"ref"`
	Ext         string `json:"ext"`
}
type FixedRespResult struct {
	Mobile      string `json:"mobile"`
	OrderId     string `json:"order_id"`
	Code        int    `json:"code"`
	ReceiveTime string `json:"receive_time"`
	Msg         string `json:"msg"`
}
type FixedResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Total  int    `json:"total"`
	Result []FixedRespResult
}

const FixedApiUrl = "https://51sms.aipaas.com/sms/sendFixedSignature"

func (sms *SmsSendService) calSig(sec string) error {
	if sec == "" {
		return errors.New("sec 不能为空")
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	tNow := time.Now().In(cstSh).Format("20060102150405")
	sms.Req.Ts = tNow
	sh := sha1.New()
	sh.Write([]byte("account=" + sms.Req.Account + "&ts=" + tNow + "&secret=" + sec))
	sms.Req.Token = fmt.Sprintf("%x", sh.Sum(nil))
	return nil
}
func (sms *SmsSendService) SendPhone(mobile, sec string) error {
	data := url.Values{}
	sms.calSig(sec)
	fmt.Println("token", sms.Req.Token)
	data.Add("account", sms.Req.Account)
	data.Add("token", sms.Req.Token)
	data.Add("ts", sms.Req.Ts)
	data.Add("mobiles", mobile)
	data.Add("content", sms.Req.Content)
	data.Add("ref", sms.Req.Ref)
	data.Add("ext", sms.Req.Ext)
	data.Add("signatureId", sms.Req.SignatureId)
	resp, err := http.PostForm(FixedApiUrl, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	readData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResp FixedResp
	err = json.Unmarshal(readData, &apiResp)
	if err != nil {
		return err
	}
	fmt.Println("---", string(readData))
	if apiResp.Code == 0 {
		return nil
	} else {
		fmt.Println("error sms api resp", string(readData), "手机号", mobile)
		return errors.New(apiResp.Msg)
	}
}
