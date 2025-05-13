package resetapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/wenlng/go-captcha-service-sdk/golang/consts"
	"github.com/wenlng/go-captcha-service-sdk/golang/resetapi/errs"
	"github.com/wenlng/go-captcha-service-sdk/golang/resetapi/request"
	"github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
	"github.com/wenlng/go-captcha-service-sdk/golang/types"
)

// Client .
type Client interface {
	GetData(ctx context.Context, id string) (*types.CaptData, error)
	CheckData(ctx context.Context, id, captchaKey, value string) (bool, error)
	CheckStatus(ctx context.Context, captchaKey string) (bool, error)
	GetStatusInfo(ctx context.Context, captchaKey string) (*types.CaptStatusInfo, error)
	DelStatusInfo(ctx context.Context, captchaKey string) (bool, error)

	UploadResource(ctx context.Context, dirname string, files []*os.File) (bool, bool, error)
	DeleteResource(ctx context.Context, id string) (bool, error)
	GetResourceList(ctx context.Context, filepath string) ([]string, error)
	GetConfig(ctx context.Context) (string, error)
	UpdateHotConfig(ctx context.Context, jsonStr string) (bool, error)
}

// client ..
type client struct {
	sdlb   *sdlb.SDLB
	config ClientConfig

	retryCount       int
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	timeout          time.Duration
}

// ClientConfig ..
type ClientConfig struct {
	BaseUrl string // http://127.0.0.1:8080
	APIKey  string // APIKeys

	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	Timeout          time.Duration

	FilterHost func(host, port string) string
}

var _ Client = (*client)(nil)

// NewHTTPClient ..
func NewHTTPClient(cnf ClientConfig, sdlb *sdlb.SDLB) (Client, error) {
	c := &client{
		sdlb:             sdlb,
		config:           cnf,
		retryCount:       3,
		retryWaitTime:    500 * time.Millisecond,
		retryMaxWaitTime: 3 * time.Second,
		timeout:          10 * time.Second,
	}

	if cnf.RetryCount > 0 {
		c.retryCount = cnf.RetryCount
	}
	if cnf.RetryWaitTime > 0 {
		c.retryWaitTime = cnf.RetryWaitTime
	}
	if cnf.RetryMaxWaitTime > 0 {
		c.retryMaxWaitTime = cnf.RetryMaxWaitTime
	}
	if cnf.Timeout > 0 {
		c.timeout = cnf.Timeout
	}

	return c, nil
}

// SelectAddress ..
func (c *client) SelectAddress() (string, error) {
	hostname, _ := os.Hostname()
	return c.SelectAddressWithKey(hostname)
}

// SelectAddressWithKey ..
func (c *client) SelectAddressWithKey(key string) (string, error) {
	if c.sdlb == nil {
		return c.config.BaseUrl, nil
	}

	inst, err := c.sdlb.Select(key)
	if err != nil {
		if c.config.BaseUrl != "" {
			return c.config.BaseUrl, nil
		}
		return "", fmt.Errorf("failed to select instance: %v", err)
	}

	if c.config.FilterHost != nil {
		return c.config.FilterHost(inst.GetHost(), inst.GetHTTPPort()), nil
	}

	url := fmt.Sprintf("http://%s", inst.GetHTTPAddress())
	return url, nil
}

// GetData ..
func (c *client) GetData(ctx context.Context, id string) (*types.CaptData, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return nil, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetQueryParam("id", id).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Get(GetDataUrlPath)

	if err != nil {
		return nil, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return nil, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return nil, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return nil, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return nil, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return nil, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return &types.CaptData{
		Id:                gjson.Get(resData.Data, "id").String(),
		CaptchaKey:        gjson.Get(resData.Data, "captcha_key").String(),
		MasterImageBase64: gjson.Get(resData.Data, "master_image_base64").String(),
		ThumbImageBase64:  gjson.Get(resData.Data, "thumb_image_base64").String(),
		MasterImageWidth:  gjson.Get(resData.Data, "master_width").Int(),
		MasterImageHeight: gjson.Get(resData.Data, "master_height").Int(),
		ThumbImageWidth:   gjson.Get(resData.Data, "thumb_width").Int(),
		ThumbImageHeight:  gjson.Get(resData.Data, "thumb_height").Int(),
		ThumbImageSize:    gjson.Get(resData.Data, "thumb_size").Int(),
		DisplayX:          gjson.Get(resData.Data, "display_x").Int(),
		DisplayY:          gjson.Get(resData.Data, "display_y").Int(),
	}, nil
}

// CheckData ..
func (c *client) CheckData(ctx context.Context, id, captchaKey, value string) (bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		SetBody(map[string]string{
			"id":         id,
			"captchaKey": captchaKey,
			"value":      value,
		}).
		Post(CheckDataUrlPath)

	if err != nil {
		return false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

  	if resData.Data == "ok" {
		return true, nil
	}
	
	return false, nil
}

// CheckStatus ..
func (c *client) CheckStatus(ctx context.Context, captchaKey string) (bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		SetQueryParam("captchaKey", captchaKey).
		Get(CheckStatusUrlPath)

	if err != nil {
		return false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return resData.Data == "ok", nil
}

// GetStatusInfo ..
func (c *client) GetStatusInfo(ctx context.Context, captchaKey string) (*types.CaptStatusInfo, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return nil, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetQueryParam("captchaKey", captchaKey).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Get(GetStatusInfoUrlPath)

	if err != nil {
		return nil, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return nil, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return nil, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return nil, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return nil, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return nil, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	ttype := gjson.Get(resData.Data, "type").Int()
	info := gjson.Get(resData.Data, "data").String()
	res := &types.CaptStatusInfo{
		Status: gjson.Get(resData.Data, "status").Int(),
		Type:   ttype,
		Info:   info,
	}

	switch ttype {
	case consts.GoCaptchaTypeClick, consts.GoCaptchaTypeClickShape:
		var data map[int]*types.ClickData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}

		res.ClickDataMaps = data
		break
	case consts.GoCaptchaTypeSlide, consts.GoCaptchaTypeDrag:
		var data map[int]*types.SlideData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}

		res.SlideDataMaps = data
		break
	case consts.GoCaptchaTypeRotate:
		var data map[int]*types.RotateData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}
		res.RotateDataMaps = data
		break
	}
	return res, nil
}

// DelStatusInfo ..
func (c *client) DelStatusInfo(ctx context.Context, captchaKey string) (bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetQueryParam("captchaKey", captchaKey).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Delete(DelStatusInfoUrlPath)

	if err != nil {
		return false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return resData.Data == "ok", nil
}

// convertFilesToFields ..
func convertFilesToFields(files []*os.File) []*resty.MultipartField {
	var fields []*resty.MultipartField

	for _, file := range files {
		_, err := file.Seek(0, 0)
		if err != nil {
			log.Printf("Warning: failed to seek file %s: %v", file.Name(), err)
			continue
		}

		field := &resty.MultipartField{
			Param:       "files",
			FileName:    file.Name(),
			Reader:      file,
			ContentType: "application/octet-stream",
		}

		fields = append(fields, field)
	}

	return fields
}

// UploadResource ..
func (c *client) UploadResource(ctx context.Context, dirname string, files []*os.File) (bool, bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"dirname": dirname,
		}).
		SetMultipartFields(convertFilesToFields(files)...).
		SetHeader("Content-Type", request.CONTENT_FORM_DATA_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Post(UploadResourceUrlPath)

	if err != nil {
		return false, false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	allDone := gjson.Get(resData.Data, "data").String()
	return true, allDone == "some-files-ok", nil
}

// DeleteResource ..
func (c *client) DeleteResource(ctx context.Context, path string) (bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetQueryParam("path", path).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Delete(DeleteResourceUrlPath)

	if err != nil {
		return false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return true, nil
}

// GetResourceList ..
func (c *client) GetResourceList(ctx context.Context, filepath string) ([]string, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return nil, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetQueryParam("path", filepath).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Get(GetResourceListUrlPath)

	if err != nil {
		return nil, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return nil, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return nil, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return nil, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return nil, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return nil, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	var list []string
	err = json.Unmarshal([]byte(resData.Data), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetConfig ..
func (c *client) GetConfig(ctx context.Context) (string, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return "", err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		Get(GetConfigUrlPath)

	if err != nil {
		return "", err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return "", errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return "", errs.ForbiddenErr
		}

		rbody := resp.Body()
		return "", fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return "", err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return "", fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return resData.Data, nil
}

// UpdateHotConfig ..
func (c *client) UpdateHotConfig(ctx context.Context, jsonStr string) (bool, error) {
	baseUrl, err := c.SelectAddress()
	if err != nil {
		return false, err
	}

	resp, err := request.GetHttpClient().
		SetBaseURL(baseUrl).
		SetRetryCount(c.retryCount).
		SetRetryWaitTime(c.retryWaitTime).
		SetRetryMaxWaitTime(c.retryMaxWaitTime).
		SetTimeout(c.timeout).
		R().
		SetContext(ctx).
		SetHeader("Content-Type", request.CONTENT_JSON_TYPE).
		SetHeader("X-API-Key", c.config.APIKey).
		SetBody(jsonStr).
		Post(UpdateHotConfigUrlPath)

	if err != nil {
		return false, err
	}

	if !request.IsHttpStatusSuccess(resp.StatusCode()) {
		if request.IsHttpStatusUnauthorized(resp.StatusCode()) {
			return false, errs.UnauthorizedErr
		} else if request.IsHttpStatusForbidden(resp.StatusCode()) {
			return false, errs.ForbiddenErr
		}

		rbody := resp.Body()
		return false, fmt.Errorf("response status: %d, result: %s", resp.StatusCode(), string(rbody))
	}

	resData, err := request.ParseResp(resp)
	if err != nil {
		return false, err
	}

	if !errs.CheckBizCodeSuccess(resData.Code) {
		return false, fmt.Errorf("code: %d, message: %v, data: %v", resData.Code, resData.Message, resData.Data)
	}

	return true, nil
}
