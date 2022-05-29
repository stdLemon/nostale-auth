package gfclient_auth

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

type TimingRange struct {
	Min int
	Max int
}

func (timing *TimingRange) Random() float64 {
	return float64(rand.Intn(timing.Max-timing.Min) + timing.Min)
}

type IdentityManager struct {
	filename string
	identity *Identity
}

type Timing struct {
	Dp TimingRange
	Df TimingRange
	Dw TimingRange
	Dc TimingRange
	D  TimingRange
}

type Identity struct {
	Timing          Timing
	Fingerprint     Fingerprint
	Installation_id string
}

func (manager *IdentityManager) loadIdentity() error {
	content, err := ioutil.ReadFile(manager.filename)

	if err != nil {
		return err
	}

	manager.identity = new(Identity)
	err = json.Unmarshal(content, manager.identity)

	if err != nil {
		return err
	}

	return nil
}

func NewIdentityManager(filename string) (IdentityManager, error) {
	manager := IdentityManager{filename: filename}
	err := manager.loadIdentity()
	return manager, err
}

func (manager *IdentityManager) Get() Identity {
	return *manager.identity
}

func (manager *IdentityManager) Update() {
	updateVector(&manager.identity.Fingerprint.Vector)
}

func (manager *IdentityManager) Save() error {
	identity_data, err := json.Marshal(manager.identity)

	if err != nil {
		return err
	}

	ioutil.WriteFile(manager.filename, identity_data, 0600)
	return nil
}
