package aligreen

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAudioAsyncScan(t *testing.T) {

	g := NewGreenClient("cn-beijing", "", "")

	imc := &AudioScanner{}
	imc.BindClient(g.c, "default")

	b, err := imc.Async([]*ContentScanData{
		&ContentScanData{
			DataId: "111",
			Url:    "",
		},
	})

	json1, _ := json.Marshal(b)
	t.Errorf("%s", json1)
	t.Error(err)

	time.Sleep(time.Duration(5) * time.Second)

	b1, err := imc.AsyncResults(b)
	json2, _ := json.Marshal(b1)
	t.Errorf("%s", json2)
	t.Error(err)

}
