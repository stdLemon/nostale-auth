package gfclient_auth

import (
	"encoding/base64"
	"time"
)

const SERVER_FILE_GAME1_FILE = "https://gameforge.com/tra/game1.js"
const VECTOR_CONTENT_LENGTH = 100
const UUID_LENGTH = 27

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
	Cookies        bool
	Mem            float64
	Con            float64
	Lang           string
	Plugins        string
	Gpu            string
	Fonts          string
	AudioC         string
	Analyser       string
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
	DP             float64
	DF             float64
	DW             float64
	DC             float64
	Creation       string
	Uuid           string
	D              float64
	OsVersion      string
	Vector         string
	UserAgent      string
	ServerTimeInMS string
	Request        *Request
}

func createFingerprint(identity_manager IdentityManager) (Fingerprint, error) {
	identity_manager.Update()

	identity := identity_manager.Get()

	fingerprint := identity.Fingerprint
	fingerprint.DP = identity.Timing.Dp.Random()
	fingerprint.DF = identity.Timing.Df.Random()
	fingerprint.DW = identity.Timing.Dw.Random()
	fingerprint.DC = identity.Timing.Dc.Random()
	fingerprint.D = identity.Timing.D.Random()

	fingerprint.Creation = time.Now().UTC().Format(time.RFC3339)
	fingerprint.Vector = base64.StdEncoding.EncodeToString([]byte(fingerprint.Vector))

	server_date, err := getServerDate()
	if err != nil {
		return Fingerprint{}, err
	}

	fingerprint.ServerTimeInMS = server_date

	return fingerprint, nil
}
