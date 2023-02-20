package tencentcloud

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"

	tc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	th "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func (p *Provider) describeRecordList(ctx context.Context, zone string) ([]libdns.Record, error) {

	list := []libdns.Record{}

	payload := map[string]any{
		"Domain": strings.Trim(zone, "."),
	}

	resp, err := p.doRequest("DescribeRecordList", payload)
	if err != nil {
		return list, err
	}

	data := DescribeRecordListResponse{}
	if err = json.Unmarshal(resp, &data); err != nil {
		return list, err
	}

	for _, record := range data.Response.RecordList {
		list = append(list, libdns.Record{
			ID:    strconv.Itoa(record.RecordId),
			Type:  record.Type,
			Name:  record.Name,
			Value: record.Value,
			TTL:   time.Duration(record.TTL) * time.Second,
		})
	}

	return list, err

}

func (p *Provider) createRecord(ctx context.Context, zone string, record libdns.Record) (string, error) {

	payload := map[string]any{
		"Domain":     strings.Trim(zone, "."),
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"RecordLine": "默认",
		"Value":      record.Value,
	}

	resp, err := p.doRequest("CreateRecord", payload)
	if err != nil {
		return "", err
	}

	data := CreateRecordResponse{}
	if err = json.Unmarshal(resp, &data); err != nil {
		return "", err
	}

	return strconv.Itoa(data.Response.RecordId), nil

}

func (p *Provider) modifyRecord(ctx context.Context, zone string, record libdns.Record) error {

	recordId, _ := strconv.Atoi(record.ID)

	payload := map[string]any{
		"Domain":     strings.Trim(zone, "."),
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"RecordLine": "默认",
		"Value":      record.Value,
		"RecordId":   recordId,
	}

	_, err := p.doRequest("ModifyRecord", payload)

	return err

}

func (p *Provider) deleteRecord(ctx context.Context, zone string, record libdns.Record) error {

	recordId, _ := strconv.Atoi(record.ID)

	payload := map[string]any{
		"Domain":   strings.Trim(zone, "."),
		"RecordId": recordId,
	}

	_, err := p.doRequest("DeleteRecord", payload)

	return err

}

func (p *Provider) doRequest(action string, payload any) ([]byte, error) {

	cpf := tp.NewClientProfile()
	cpf.HttpProfile.RootDomain = "tencentcloudapi.com"

	cred := tc.NewCredential(p.SecretId, p.SecretKey)
	client := tc.NewCommonClient(cred, "", cpf)

	request := th.NewCommonRequest("dnspod", "2021-03-23", action)
	request.SetActionParameters(payload)

	response := th.NewCommonResponse()

	if err := client.Send(request, response); err != nil {
		return nil, err
	}

	return response.GetBody(), nil

}
