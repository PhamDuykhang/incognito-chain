package ringct

type BulletProof struct {
	V []privacy.Key // 1 * 32           // extra 1 byte for length

	// 4
	A  crypto.Key // 1 * 32
	S  crypto.Key // 1 * 32
	T1 crypto.Key // 1 * 32
	T2 crypto.Key // 1 * 32

	// final 2/5
	taux crypto.Key // 1 * 32
	mu   crypto.Key // 1 * 32

	// 2*6
	L []crypto.Key // 6 * 32  // space requirements while serializing, extra 1 byte for length
	R []crypto.Key // 6 * 32 // space requirements while serializing, extra 1 byte for length

	// final 3/5
	a crypto.Key // 1 * 32
	b crypto.Key // 1 * 32
	t crypto.Key // 1 * 32
}
