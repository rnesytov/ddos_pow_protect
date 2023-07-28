import scrypt

# Took from litecoin scrypt parameters
scrypt_N = 1 << 10
scrypt_r = 1
scrypt_p = 1

from hashlib import sha256

def generate_proof_of_work(challenge: str, difficulty: int) -> int:
    nonce = 0
    while True:
        hash = scrypt.hash(
            challenge + str(nonce), challenge, N=scrypt_N, p=scrypt_p, r=scrypt_r
        )
        if hash.hex().startswith("0" * difficulty):
            return nonce
        nonce += 1


def verify_proof_of_work(challenge: str, nonce: int, difficulty: int) -> bool:
    hash = scrypt.hash(
        challenge + str(nonce), challenge, N=scrypt_N, p=scrypt_p, r=scrypt_r
    )
    return hash.hex().startswith("0" * difficulty)
