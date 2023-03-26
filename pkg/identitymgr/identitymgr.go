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

func (m *Manager) loadIdentity() error {
	b, err := os.ReadFile(m.filename)

	if err != nil {
		return err
	}

	m.identity = new(Identity)
	err = json.Unmarshal(b, m.identity)

	if err != nil {
		return err
	}

	return nil
}

func New(filename string) (Manager, error) {
	m := Manager{filename: filename}
	err := m.loadIdentity()
	return m, err
}

func (m *Manager) Get() Identity {
	return *m.identity
}

func (m *Manager) Save() error {
	const USER_RW_PERMS = 0600

	b, err := json.Marshal(m.identity)

	if err != nil {
		return err
	}

	os.WriteFile(m.filename, b, USER_RW_PERMS)
	return nil
}

func (m *Manager) CreateFingerprint() (blackbox.Fingerprint, error) {
	i := m.identity

	f := i.Fingerprint
	blackbox.UpdateVector(&f.Vector)

	f.D = i.Timing.Random()

	f.Creation = time.Now().UTC().Format(time.RFC3339)
	f.Vector = base64.StdEncoding.EncodeToString([]byte(f.Vector))

	date, err := blackbox.GetServerDate()
	if err != nil {
		return blackbox.Fingerprint{}, err
	}

	f.ServerTimeInMS = date

	return f, nil
}

func (m *Manager) NewBlackbox(r *blackbox.Request) (blackbox.Blackbox, error) {
	f, err := m.CreateFingerprint()
	if err != nil {
		return "", err
	}

	f.Request = r
	return blackbox.New(&f)
}

func (m *Manager) NewEncryptedBlackbox(gsId, accountId string) ([]byte, error) {
	i := strings.LastIndexByte(gsId, '-')
	session := gsId[:i]

	const featureMin = 1
	const featureMax = 0xFFFFFFFE
	feature := float64(rand.Intn(featureMax-featureMin) + featureMin)
	request := blackbox.Request{Features: []float64{feature}, Installation: m.identity.InstallationId, Session: session}

	blackbox, err := m.NewBlackbox(&request)

	if err != nil {
		return nil, err
	}

	return blackbox.Encrypt(gsId, accountId), nil
}
