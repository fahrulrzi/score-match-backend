package entity

import "time"

type Customer struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Job       string    `json:"job"`
	Income    int       `json:"income"`
	Age       int       `json:"age"`
	Score     int       `json:"score"`
	Status    string    `json:"status"`
	Describe  string    `json:"describe"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerRequest struct {
	Username      string `json:"username"`
	Job           string `json:"job"`
	Income        int    `json:"income"`
	Installment   int    `json:"installment"`
	Age           int    `json:"age"`
	MaritalStatus string `json:"marital_status"`
	LengthOfWork  string `json:"length_of_work"`
	Purpose       string `json:"purpose"`
	Collateral    string `json:"collateral"`
}

type CustomerScoreResponse struct {
	Score    int    `json:"score"`
	Status   string `json:"status"`
	Describe string `json:"describe"`
}

var MaritalStatusScore = map[string]int{
	"Lajang":   100,
	"Menikah":  80,
	"Bercerai": 70,
}

type DBRScoreRange struct {
	Min   float64
	Max   float64
	Score int
}

var DBRScoreRanges = []DBRScoreRange{
	{Min: 0, Max: 10, Score: 100},
	{Min: 11, Max: 20, Score: 80},
	{Min: 21, Max: 30, Score: 70},
	{Min: 31, Max: 40, Score: 60},
	{Min: 41, Max: 1000, Score: 10},
}

var JobScore = map[string]int{
	"Pegawai Tetap PNS,BUMN": 100,
	"Pegawai Swasta":         90,
	"Pegawai Kontrak":        70,
	"Freelance":              60,
	"Serabutan":              40,
	"Tidak Bekerja":          0,
}

var LengthOfWorkScore = map[string]int{
	">5 tahun":  100,
	"3-5 tahun": 80,
	"1-3 tahun": 70,
	"<1 tahun":  30,
}

var Purpose = map[string]int{
	"Modal Usaha": 100,
	"Konsumtif":   70,
}

var Collateral = map[string]int{
	"Sertifikat Rumah": 100,
	"BPKB Mobil":       80,
	"BPKB Motor":       70,
	"Tidak ada agunan": 0,
}
