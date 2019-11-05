package watch

import (
	"regexp"
)

/*
{
  "kind": "storage#object",
  "id": "zedge-test-chartmuseum/evtail-server-v1.1.19-23-g235d9b4.tgz/1572813538165970",
  "selfLink": "https://www.googleapis.com/storage/v1/b/zedge-test-chartmuseum/o/evtail-server-v1.1.19-23-g235d9b4.tgz",
  "name": "evtail-server-v1.1.19-23-g235d9b4.tgz",
  "bucket": "zedge-test-chartmuseum",
  "generation": "1572813538165970",
  "metageneration": "1",
  "contentType": "application/x-gzip",
  "timeCreated": "2019-11-03T20:38:58.165Z",
  "updated": "2019-11-03T20:38:58.165Z",
  "storageClass": "REGIONAL",
  "timeStorageClassUpdated": "2019-11-03T20:38:58.165Z",
  "size": "2586",
  "md5Hash": "ojxKFCng3E2tOT+bO/iK0A==",
  "mediaLink": "https://www.googleapis.com/download/storage/v1/b/zedge-test-chartmuseum/o/evtail-server-v1.1.19-23-g235d9b4.tgz?generation=1572813538165970&alt=media",
  "crc32c": "XBtpAg==",
  "etag": "CNKJ9oHzzuUCEAE="
}
*/
type ChartMessage struct {
	Kind                    string `json:"kind"`
	ID                      string `json:"id"`
	SelfLink                string `json:"self_link"`
	Name                    string `json:"name"`
	Bucket                  string `json:"bucket"`
	Generation              string `json:"generation"`
	MetaGeneration          string `json:"metageneration"`
	ContentType             string `json:"contentType"`
	TimeCreated             string `json:"timeCreated"`
	Updated                 string `json:"updated"`
	StorageClass            string `json:"storageClass"`
	TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
	Size                    string `json:"size"`
	MD5Hash                 string `json:"md5Hash"` // base64 encoded
	MediaLink               string `json:"mediaLink"`
	CRC32c                  string `json:"crc32c"` // base64 encoded
	ETag                    string `json:"etag"`
}

var nameVersionRegex = regexp.MustCompile(`^(.*)-(v?\d+\..*?)\.tgz$`)

func (m ChartMessage) GetChartNameAndVersion() (string, string) {
	matches := nameVersionRegex.FindStringSubmatch(m.Name)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}
