from pow import generate_proof_of_work
from server import get_challenge
from time import time


def main():
    n = 500
    d = 3
    total_time = 0.0

    print(f"{n=}, {d=}")

    for _ in range(n):
        t0 = time()
        generate_proof_of_work(get_challenge(32), d)
        total_time += time() - t0

    print(f"Total={total_time}, average={total_time/n}")


if __name__ == "__main__":
    main()
