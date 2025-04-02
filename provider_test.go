package tencentcloud

import (
	"context"
	"os"
	"testing"

	"github.com/libdns/libdns"
)

var provider = &Provider{
	SecretId:  os.Getenv("TC_SECRET_ID"),
	SecretKey: os.Getenv("TC_SECRET_KEY"),
}

var (
	zone  = os.Getenv("TC_ZONE")
	name  = os.Getenv("TC_NAME")
	value = os.Getenv("TC_VALUE")
)

func TestSetRecords(t *testing.T) {
	_, err := provider.SetRecords(context.Background(), zone, []libdns.Record{
		{
			Type:  "A",
			Name:  name,
			Value: value,
		},
	})
	if err != nil {
		t.Fatalf("SetRecords: %v", err)
	}
}

func TestGetRecords(t *testing.T) {
	records, err := provider.GetRecords(context.Background(), zone)
	if err != nil {
		t.Fatalf("GetRecords: %v", err)
	}
	for _, record := range records {
		t.Logf("RecordId: %s, RecordType: %s, Name: %s, Value: %s, TTL: %d",
			record.ID, record.Type, record.Name, record.Value, int(record.TTL.Seconds()))
	}
}
