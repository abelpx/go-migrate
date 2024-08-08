package mysql

type metadata struct {
	Name          string
	Type          string
	Length        int
	Decimals      int
	Nullable      bool
	Unique        indexName
	Index         indexName
	Modify        bool
	IndexName     string
	Primary       bool
	AutoIncrement bool
	unsigned      bool
	Collate       string
	Default       interface{}
	Comment       string
	Custom        string
	Foreign       *foreign
	TableComment  string
}
