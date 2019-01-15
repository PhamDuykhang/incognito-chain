package zkp

import (
	"fmt"
	"github.com/ninjadotorg/constant/privacy"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

//TestInnerProduct test inner product calculation
func TestInnerProduct(t *testing.T) {
	n := 2
	a := make([]*big.Int, n)
	b := make([]*big.Int, n)

	for i:=0; i<n; i++{
		a[i]= big.NewInt(10)
		b[i]= big.NewInt(20)
	}

	c, _ := innerProduct(a, b)
	assert.Equal(t, big.NewInt(400), c)

	bytes := privacy.RandBytes(33)

	num1 := new(big.Int).SetBytes(bytes)
	num1Inverse := new(big.Int).ModInverse(num1, privacy.Curve.Params().N)

	num2 := new(big.Int).SetBytes(bytes)
	num2 = num2.Mod(num2, privacy.Curve.Params().N)
	num2Inverse := new(big.Int).ModInverse(num2, privacy.Curve.Params().N)

	assert.Equal(t, num1Inverse, num2Inverse)
}

func TestProve(t *testing.T){
	wit := new(InnerProductWitness)
	n := 64
	wit.a = make([]*big.Int, n)
	wit.b = make([]*big.Int, n)
	for i := range wit.a{
		wit.a[i] = big.NewInt(10)
		wit.b[i] = big.NewInt(10)
	}

	wit.u = new(privacy.EllipticPoint)
	wit.u.Randomize()

	wit.p = new(privacy.EllipticPoint).Zero()
	c, err := innerProduct(wit.a, wit.b)
	if err != nil{
		fmt.Printf("Err: %v\n", err)
	}

	for i := range wit.a{
		wit.p = wit.p.Add(AggParam.G[i].ScalarMult(wit.a[i]))
		wit.p = wit.p.Add(AggParam.H[i].ScalarMult(wit.b[i]))
	}
	wit.p = wit.p.Add(wit.u.ScalarMult(c))

	proof, err:= wit.Prove()
	if err != nil{
		fmt.Printf("Err: %v\n", err)
	}

	fmt.Printf("Proving done!!!")
	res := proof.Verify()

	assert.Equal(t, true, res)
}
