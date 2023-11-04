package blackbox

const (
	SERVER_FILE_GAME1_FILE = "https://gameforge.com/tra/game1.js"
	VECTOR_CONTENT_LENGTH  = 100
	UUID_LENGTH            = 27
)

type Request struct {
	Features     []float64 `json:"features"`
	Installation string    `json:"installation"`
	Session      string    `json:"session"`
}

type Fingerprint struct {
	V              float64  `json:"v"`
	Tz             string   `json:"tz"`
	Dnt            bool     `json:"dnt"`
	Product        string   `json:"product"`
	OsType         string   `json:"osType"`
	App            string   `json:"app"`
	Vendor         string   `json:"vendor"`
	Mem            float64  `json:"mem"`
	Con            float64  `json:"con"`
	Lang           string   `json:"lang"`
	Plugins        string   `json:"plugins"`
	Gpu            string   `json:"gpu"`
	Fonts          string   `json:"fonts"`
	AudioC         string   `json:"audioC"`
	Width          float64  `json:"width"`
	Height         float64  `json:"height"`
	Depth          float64  `json:"depth"`
	Video          string   `json:"video"`
	Audio          string   `json:"audio"`
	Media          string   `json:"media"`
	Permissions    string   `json:"permissions"`
	AudioFP        float64  `json:"audioFP"`
	WebglFP        string   `json:"webglFP"`
	CanvasFP       float64  `json:"canvasFP"`
	Creation       string   `json:"creation"`
	Uuid           string   `json:"uuid"`
	D              float64  `json:"d"`
	OsVersion      string   `json:"osVersion"`
	Vector         string   `json:"vector"`
	UserAgent      string   `json:"userAgent"`
	ServerTimeInMS string   `json:"serverTimeInMS"`
	Request        *Request `json:"request"`
}
