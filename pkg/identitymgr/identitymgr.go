package identitymgr

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/stdLemon/nostale-auth/pkg/blackbox"
)

type Manager struct {
	filename string
	identity *Identity
}

func New(filename string) (Manager, error) {
	mgr := Manager{filename: filename}

	content, err := os.ReadFile(mgr.filename)
	if err != nil {
		return mgr, err
	}

	mgr.identity = new(Identity)
	if err := json.Unmarshal(content, mgr.identity); err != nil {
		return mgr, err
	}

	return mgr, nil
}

func (m Manager) Get() Identity {
	return *m.identity
}

func (m Manager) Save() error {
	const USER_RW_PERMS = 0600

	jsonIdentity, err := json.Marshal(m.identity)
	if err != nil {
		return err
	}

	return os.WriteFile(m.filename, jsonIdentity, USER_RW_PERMS)
}

func (m Manager) createFingerprint(r *blackbox.Request) (blackbox.Fingerprint, error) {
	var (
		err         error
		fingerprint = m.identity.Fingerprint
	)

	fingerprint.ServerTimeInMS, err = blackbox.GetServerDate()
	if err != nil {
		return blackbox.Fingerprint{}, err
	}

	m.identity.Fingerprint.Vector = blackbox.UpdateVector(fingerprint.Vector)

	fingerprint.Request = r
	fingerprint.D = m.identity.Timing.Random()
	fingerprint.Creation = time.Now().UTC().Format(time.RFC3339)
	fingerprint.Vector = base64.StdEncoding.EncodeToString([]byte(fingerprint.Vector))
	return fingerprint, nil
}

func (m Manager) NewBlackbox(request *blackbox.Request) (blackbox.Blackbox, error) {
	fingerprint, err := m.createFingerprint(request)
	if err != nil {
		return "", err
	}

	return blackbox.New(&fingerprint)
}

func (m Manager) NewEncryptedBlackbox(gsId, accountId string) ([]byte, error) {
	const (
		featureMin = 1
		featureMax = 0xFFFFFFFE
	)

	var (
		session = gsId[:strings.LastIndexByte(gsId, '-')]
		feature = float64(rand.Intn(featureMax-featureMin) + featureMin)
		request = blackbox.Request{Features: []float64{feature}, Installation: m.identity.InstallationId, Session: session}
	)

	blackbox, err := m.NewBlackbox(&request)
	if err != nil {
		return nil, err
	}

	return blackbox.Encrypt(gsId, accountId), nil
}
