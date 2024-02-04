package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"signer/config"
	"signer/repo"
	"strings"
)

type TestService interface {
	Sign(ctx context.Context, username string, testQA string) (string, error)
	CheckSignature(ctx context.Context, username string, signature string) (string, string, error)
}

type testService struct {
	config *config.Config
	tr     *repo.TestRepo
	qar    *repo.QA
}

func NewTestService(cfg *config.Config, tr *repo.TestRepo, qar *repo.QA) TestService {
	return &testService{config: cfg, tr: tr, qar: qar}
}

func (ts *testService) Sign(ctx context.Context, username string, testQA string) (string, error) {
	sig := sha256.Sum256([]byte(testQA))
	sigHex := hex.EncodeToString(sig[:])
	id, err := ts.tr.SignTest(username, sigHex)
	if err != nil {
		return "", err
	}
	QAs := strings.Split(testQA, ";")
	for _, qa := range QAs {
		spl := strings.Split(qa, ":")
		if err := ts.qar.InsertQAs(strings.TrimSpace(spl[0]), strings.TrimSpace(spl[1]), id); err != nil {
			return "", err
		}
	}

	return sigHex, nil
}

func (ts *testService) CheckSignature(ctx context.Context, username string, signature string) (string, string, error) {
	sig, err := ts.tr.GetSignature(username)
	if err != nil {
		return "", "", err
	}
	if sig.Signature == signature {
		return "", "", errors.New("sinature not matching")
	}

	QAs, err := ts.qar.GetQAs(sig.SignatureID)
	if err != nil {
		return "", "", err
	}

	return sig.Timestamp.String(), QAs, nil
}
