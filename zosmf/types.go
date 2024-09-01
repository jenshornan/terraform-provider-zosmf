package zosmf

type CreateDataset struct {
	Volser    string `json:"volser"`
	Unit      string `json:"unit"`
	Dsorg     string `json:"dsorg"`
	Alcunit   string `json:"alcunit"`
	Primary   int    `json:"primary"`
	Secondary int    `json:"secondary"`
	Dirblk    int    `json:"dirblk"`
	Avgblk    int    `json:"avgblk"`
	Recfm     string `json:"recfm"`
	Blksize   int    `json:"blksize"`
	Lrecl     int    `json:"lrecl"`
	Storclass string `json:"storclass"`
	Mgntclass string `json:"mgntclass"`
	Dataclass string `json:"dataclass"`
	Dsntype   string `json:"dsntype"`
}

type Dataset struct {
	Content string
}
