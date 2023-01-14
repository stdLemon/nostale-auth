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
	V              float64
	Tz             string
	Dnt            bool
	Product        string
	OsType         string
	App            string
	Vendor         string
	Mem            float64
	Con            float64
	Lang           string
	Plugins        string
	Gpu            string
	Fonts          string
	AudioC         string
	Width          float64
	Height         float64
	Depth          float64
	LStore         bool
	SStore         bool
	Video          string
	Audio          string
	Media          string
	Permissions    string
	AudioFP        float64
	WebglFP        string
	CanvasFP       float64
	Creation       string
	Uuid           string
	D              float64
	OsVersion      string
	Vector         string
	UserAgent      string
	ServerTimeInMS string
	Request        *Request
}
