package tencentcloud

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	endpoint     = "https://dnspod.tencentcloudapi.com"
	reqJson      = `{"RecordType":"","Domain":"","RecordLine":"默认","SubDomain":"","Value":"","RecordId":0}`
	reqJson_find = `{"RecordType":"","Domain":"","Subdomain":""}`

	DescribeRecordList = "DescribeRecordList"
	CreateRecord       = "CreateRecord"
	ModifyRecord       = "ModifyRecord"
	DeleteRecord       = "DeleteRecord"
)

var sOption = sjson.Options{Optimistic: true, ReplaceInPlace: true}

func (p *Provider) listRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	domain := strings.TrimSuffix(zone, ".")
	payload, err := sjson.Set("", "Domain", domain)
	if err != nil {
		return nil, err
	}

	resp, err := p.sendRequest(ctx, DescribeRecordList, payload)
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(resp, "Response.RecordList")
	if !result.IsArray() {
		return nil, ErrNotValid
	}

	list := make([]libdns.Record, 0, result.Get("#").Int())
	result.ForEach(func(_, v gjson.Result) bool {
		list = append(list, libdns.Record{
			ID:    v.Get("RecordId").String(),
			Type:  v.Get("Type").String(),
			Name:  v.Get("Name").String(),
			Value: v.Get("Value").String(),
			TTL:   time.Duration(v.Get("TTL").Int()) * time.Second,
		})
		return true
	})

	return list, nil
}

func (p *Provider) createRecord(ctx context.Context, zone string, record libdns.Record) error {
	domain := strings.TrimSuffix(zone, ".")

	payload, _ := sjson.SetOptions(reqJson, "Domain", domain, &sOption)
	payload, _ = sjson.SetOptions(payload, "SubDomain", record.Name, &sOption)
	payload, _ = sjson.SetOptions(payload, "RecordType", record.Type, &sOption)
	payload, _ = sjson.SetOptions(payload, "Value", record.Value, &sOption)
	payload, _ = sjson.Delete(payload, "RecordId")

	resp, err := p.sendRequest(ctx, CreateRecord, payload)
	if err != nil {
		return err
	}

	result := gjson.GetBytes(resp, "Response.RecordId")
	if !result.Exists() {
		return ErrNotValid
	}

	return nil
}

func (p *Provider) modifyRecord(ctx context.Context, zone string, record libdns.Record) error {
	domain := strings.TrimSuffix(zone, ".")

	payload, _ := sjson.SetOptions(reqJson, "Domain", domain, &sOption)
	payload, _ = sjson.SetOptions(payload, "SubDomain", record.Name, &sOption)
	payload, _ = sjson.SetOptions(payload, "RecordType", record.Type, &sOption)
	payload, _ = sjson.SetOptions(payload, "Value", record.Value, &sOption)
	payload, _ = sjson.SetOptions(payload, "RecordId", p.id, &sOption)

	_, err := p.sendRequest(ctx, ModifyRecord, payload)
	return err
}

func (p *Provider) deleteRecord(ctx context.Context, zone string, record libdns.Record) error {
	domain := strings.TrimSuffix(zone, ".")

	payload, _ := sjson.Set("", "Domain", domain)
	payload, _ = sjson.Set(payload, "RecordId", record.ID)

	_, err := p.sendRequest(ctx, DeleteRecord, payload)
	return err
}

func (p *Provider) findRecord(ctx context.Context, zone string, record libdns.Record) error {
	domain := strings.TrimSuffix(zone, ".")

	payload, _ := sjson.SetOptions(reqJson_find, "Domain", domain, &sOption)
	payload, _ = sjson.SetOptions(payload, "RecordType", record.Type, &sOption)
	payload, _ = sjson.SetOptions(payload, "Subdomain", record.Name, &sOption)

	resp, err := p.sendRequest(ctx, DescribeRecordList, payload)
	if err != nil {
		return err
	}

	result := gjson.GetBytes(resp, "Response.RecordList.0.RecordId")
	if !result.Exists() {
		return ErrRecordNotFound
	}

	p.id = result.Uint()
	return nil
}

func (p *Provider) sendRequest(ctx context.Context, action string, data string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-Version", "2021-03-23")

	SignRequest(p.SecretId, p.SecretKey, req, action, data)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
