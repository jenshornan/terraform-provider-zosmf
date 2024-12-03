package zosmf

type DatasetAttribute struct {
	Volser    string `json:"volser,omitempty"`
	Unit      string `json:"unit,omitempty"`
	Dsorg     string `json:"dsorg,omitempty"`
	Alcunit   string `json:"alcunit,omitempty"`
	Primary   int    `json:"primary,omitempty"`
	Secondary int    `json:"secondary,omitempty"`
	Dirblk    int    `json:"dirblk,omitempty"`
	Avgblk    int    `json:"avgblk,omitempty"`
	Recfm     string `json:"recfm,omitempty"`
	Blksize   int    `json:"blksize,omitempty"`
	Lrecl     int    `json:"lrecl,omitempty"`
	Storclass string `json:"storclass,omitempty"`
	Mgntclass string `json:"mgntclass,omitempty"`
	Dataclass string `json:"dataclass,omitempty"`
	Dsntype   string `json:"dsntype,omitempty"`
}

type Dataset struct {
	Content string
}
