package pow

import (
	"testing"
)

func TestGene(t *testing.T) {
	pow := New(NewDefaultScryptConf())
	challenge := pow.GetChallenge(32)

	nonce, err := pow.Solve(challenge, 2)
	if err != nil {
		t.Fatal(err)
	}
	valid, err := pow.Verify(challenge, 2, nonce)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("invalid nonce")
	}
}

func BenchmarkSolve2(b *testing.B) {
	pow := New(NewDefaultScryptConf())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		challenge := pow.GetChallenge(32)
		pow.Solve(challenge, 2)
	}
}

func BenchmarkSolve3(b *testing.B) {
	pow := New(NewDefaultScryptConf())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		challenge := pow.GetChallenge(32)
		pow.Solve(challenge, 3)
	}
}

func BenchmarkVerify(b *testing.B) {
	pow := New(NewDefaultScryptConf())
	challenge := pow.GetChallenge(32)
	nonce, err := pow.Solve(challenge, 3)
	if err != nil {
		b.Fatal(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.Verify(challenge, 3, nonce)
	}
}
