package gomaxcompute

type resultMeta struct {
	DownloadID  string `json:"DownloadID"`
	RecordCount int64  `json:"RecordCount"`
	Schema      struct {
		IsVirtualView bool `json:"IsVirtualView"`
		Columns       []struct {
			Comment  string `json:"comment"`
			Name     string `json:"name"`
			Nullable bool   `json:"nullable"`
			Type     string `json:"type"`
		} `json:"columns"`
	} `json:"Schema"`
	// not parsed: PartitionKeys []string `json:"partitionKeys"`
	Status string `json:"Status"`
}
