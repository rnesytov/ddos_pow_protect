package pow

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
)

const KEY_SIZE = 32

type scryptConf struct {
	N int
	p int
	r int
}

type PoW struct {
	sConf *scryptConf
}

func New(conf *scryptConf) *PoW {
	return &PoW{
		sConf: conf,
	}
}

func NewDefaultScryptConf() *scryptConf {
	return &scryptConf{
		N: 1 << 10,
		p: 1,
		r: 1,
	}
}

func (p *PoW) Verify(challenge []byte, difficulty uint8, nonce uint64) (bool, error) {
	toHash := make([]byte, len(challenge)+8)
	copy(toHash, challenge)
	binary.BigEndian.PutUint64(toHash, nonce)
	key, err := scrypt.Key(toHash, challenge, p.sConf.N, p.sConf.r, p.sConf.p, KEY_SIZE)
	if err != nil {
		return false, err
	}
	return countLeadingZeroes(key) >= int(difficulty), nil
}

func (p *PoW) GetChallenge(challengeLen uint) []byte {
	data := make([]byte, challengeLen)
	rand.Read(data)
	return data
}

func (p *PoW) Solve(challenge []byte, difficulty uint8) (uint64, error) {
	nonce := uint64(0)
	for {
		valid, err := p.Verify(challenge, difficulty, nonce)
		if err != nil {
			return 0, err
		}
		if valid {
			break
		}
		nonce++
	}
	return nonce, nil
}

func countLeadingZeroes(hash []byte) int {
	count := 0
	hexed := hex.EncodeToString(hash)
	for _, r := range hexed {
		if r == '0' {
			count += 1
		} else {
			break
		}
	}
	return count
}
