package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions sdk.Dec, phase, startPhaseBlock uint64) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
		Phase:            phase,
		StartPhaseBlock:  startPhaseBlock,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		0,
		0,
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(13, 2),
	)
}

// validate minter
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// PhaseInflationRate returns the inflation rate by phase.
func (m Minter) PhaseInflationRate(phase uint64) sdk.Dec {
	switch {
	case phase > 15:
		return sdk.ZeroDec()

	case phase == 1:
		return sdk.NewDecWithPrec(25, 2)

	case phase == 2:
		return sdk.NewDecWithPrec(15, 2)

	case phase == 3, phase == 4:
		return sdk.NewDecWithPrec(10, 2)

	case phase == 5, phase == 6, phase == 7, phase == 8, phase == 9, phase == 10:
		return sdk.NewDecWithPrec(5, 2)

	case phase == 11, phase == 12, phase == 13, phase == 14, phase == 15:
		return sdk.NewDecWithPrec(2, 2)

	default:
		return sdk.ZeroDec()
	}
}

// NextPhase returns the new phase.
func (m Minter) NextPhase(params Params, currentBlock uint64) uint64 {
	nonePhase := m.Phase == 0
	if nonePhase {
		return 1
	}

	blockNewPhase := m.StartPhaseBlock + params.BlocksPerYear
	if blockNewPhase > currentBlock {
		return m.Phase
	}

	return m.Phase + 1
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
