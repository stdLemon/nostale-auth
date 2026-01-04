package blackbox

const (
	_gfAuthScriptUrl     = "https://gameforge.com/tra/game1.js"
	_vectorContentLength = 100
	_uuidLength          = 27
)

type Request struct {
	Features     []float64 `json:"features"`
	Installation string    `json:"installation"`
	Session      string    `json:"session"`
}

type Fingerprint struct {
	SchemaVersion           float64  `json:"schemaVersion"`
	TimeZone                string   `json:"timeZone"`
	OSName                  string   `json:"osName"`
	BrowserName             string   `json:"browserName"`
	BrowserVendor           string   `json:"browserVendor"`
	DeviceMemoryGB          float64  `json:"deviceMemoryGb"`
	HardwareConcurrency     float64  `json:"hardwareConcurrency"`
	Languages               string   `json:"languages"`
	PluginsHash             string   `json:"pluginsHash"`
	WebGLVendorRenderer     string   `json:"webglVendorRenderer"`
	FontProbeHash           string   `json:"fontProbeHash"`
	AudioContextHash        string   `json:"audioContextHash"`
	ScreenAvailWidth        float64  `json:"screenAvailWidth"`
	ScreenAvailHeight       float64  `json:"screenAvailHeight"`
	VideoCodecSupportHash   string   `json:"videoCodecSupportHash"`
	AudioCodecSupportHash   string   `json:"audioCodecSupportHash"`
	MediaDeviceKindsHash    string   `json:"mediaDeviceKindsHash"`
	PermissionStatesHash    string   `json:"permissionStatesHash"`
	OfflineAudioFingerprint float64  `json:"offlineAudioFingerprint"`
	WebGLPixelHash          string   `json:"webglPixelHash"`
	CanvasFingerprintHash   float64  `json:"canvasFingerprintHash"`
	GeneratedAtISO          string   `json:"generatedAtIso"`
	ClientID                string   `json:"clientId"`
	CollectionDurationMs    float64  `json:"collectionDurationMs"`
	OSVersion               string   `json:"osVersion"`
	VecSignatureBase64      string   `json:"vecSignatureBase64"`
	UserAgent               string   `json:"userAgent"`
	ServerDateISO           string   `json:"serverDateIso"`
	ExtraPayload            *Request `json:"extraPayload"`
	AutomationFlags         float64  `json:"automationFlags"`
}
