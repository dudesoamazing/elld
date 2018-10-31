// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package blakimoto

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ellcrys/elld/params"
	"github.com/ellcrys/elld/types"
	"github.com/ellcrys/elld/types/core"
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// VerifyHeader checks whether a header
// conforms to the consensus rules
func (b *Blakimoto) VerifyHeader(header, parent types.Header, seal bool) error {

	// Ensure that the header's extra-data
	// section is of a reasonable size
	if uint64(len(header.GetExtra())) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.GetExtra()),
			params.MaximumExtraDataSize)
	}

	// Verify the header's timestamp
	if time.Unix(header.GetTimestamp(), 0).After(time.Now().
		Add(params.AllowedFutureBlockTime)) {
		return ErrFutureBlock
	}

	if header.GetTimestamp() <= parent.GetTimestamp() {
		return errZeroBlockTime
	}

	// Verify the block's difficulty based on
	// it's timestamp and parent's difficulty
	expected := b.CalcDifficulty(header, parent)
	if expected.Cmp(header.GetDifficulty()) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v",
			header.GetDifficulty(), expected)
	}

	// Verify that the total difficulty is
	// parent total difficulty + header total
	// difficulty
	expectedTd := new(big.Int).Add(parent.GetTotalDifficulty(), header.GetDifficulty())
	if headerTd := header.GetTotalDifficulty(); headerTd.Cmp(expectedTd) != 0 {
		return fmt.Errorf("invalid total difficulty: have %v, want %v",
			headerTd, expectedTd)
	}

	// Verify that the block number is
	// parent's +1
	if diff := header.GetNumber() - parent.GetNumber(); diff != 1 {
		return ErrInvalidNumber
	}

	// Verify the engine specific seal
	// securing the block
	if seal {
		if err := b.VerifySeal(header); err != nil {
			return err
		}
	}

	return nil
}

// CalcDifficulty is the difficulty adjustment
// algorithm. It returns the difficulty that a
// new block should have when created at time
// given the parent block's time and difficulty.
func (b *Blakimoto) CalcDifficulty(blockHeader types.Header, parent types.Header) *big.Int {

	if sameDiffEpoch(parent.GetNumber(), blockHeader.GetNumber()) {
		return parent.GetDifficulty()
	}

	return CalcDifficulty(blockHeader, parent)
}

// CalcDifficulty is the difficulty adjustment
// algorithm. It returns the difficulty that a new
// block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(blockHeader types.Header, parent types.Header) *big.Int {
	return calcDifficultyInception(uint64(blockHeader.GetTimestamp()),
		parent)
}

// Some weird constants to avoid constant memory
// allocs for them.
var (
	expDiffPeriod = big.NewInt(3)
	big0          = big.NewInt(0)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big100F       = big.NewFloat(100)
)

func sameDiffEpoch(parentBlockNumber, blockNumber uint64) bool {
	parentEpoch := int(parentBlockNumber / params.DifficultyEpoch)
	blockEpoch := int(blockNumber / params.DifficultyEpoch)
	return parentEpoch == blockEpoch
}

func calcDifficultyInception(time uint64, parent types.Header) *big.Int {

	diff := new(big.Int)
	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(new(big.Int).
		SetInt64(parent.GetTimestamp()))

	// Define the value to adjust difficulty by
	// when the block time is within or above the duration limit.
	adjust := new(big.Int).Div(parent.GetDifficulty(), params.DifficultyBoundDivisor)

	// Calculate the time difference between the
	// new block time and parent block time. We'll
	// use this to determine whether on not to increase
	// or decrease difficulty
	blockTimeDiff := bigTime.Sub(bigTime, bigParentTime)

	// When block time difference is within the expected
	// block duration limit, we increasse difficulty
	if blockTimeDiff.Cmp(params.DurationLimit) < 0 {
		diff.Add(parent.GetDifficulty(), adjust)
	}

	// When block time difference is equal or greater than
	// the expected block duration limit, we decrease difficulty
	if blockTimeDiff.Cmp(params.DurationLimit) >= 0 {

		// We need to determine the percentage increase of the
		// new block time compared to the duration limit
		durLimitF := new(big.Float).SetInt(params.DurationLimit)
		blockTimeF := new(big.Float).SetInt(bigTime)
		timeDiff := new(big.Float).Sub(blockTimeF, durLimitF)
		pctDiff := new(big.Float).Mul(new(big.Float).Quo(timeDiff, durLimitF), big100F)

		// If the percentage difference is below the allowed mimimum
		// reset to the minimum
		if pctDiff.Cmp(params.MinimumDurationIncrease) < 0 {
			pctDiff = new(big.Float).Set(params.MinimumDurationIncrease)
		}

		// Calculate the new adjustment based on time difference percentage
		pctDiff = pctDiff.Quo(pctDiff, big100F)
		timeDiffAdjust, _ := new(big.Float).Mul(pctDiff, new(big.Float).SetInt(adjust)).Int(nil)
		diff.Sub(parent.GetDifficulty(), timeDiffAdjust)
	}

	// Ensure difficulty does not go below the required minimum
	if diff.Cmp(params.MinimumDifficulty) < 0 {
		diff.Set(params.MinimumDifficulty)
	}

	return diff
}

// VerifySeal checks whether the given
// block satisfies the PoW difficulty
// requirements.
func (b *Blakimoto) VerifySeal(header types.Header) error {

	// If we're running a fake PoW, accept any seal as valid
	if b.config.PowMode == ModeTest {
		time.Sleep(b.fakeDelay)
		// if b.fakeFail == header.Number {
		// 	return errInvalidPoW
		// }
		return nil
	}

	// Ensure that we have a valid difficulty for the block
	if header.GetDifficulty().Sign() <= 0 {
		return errInvalidDifficulty
	}

	// Recompute the digest and PoW value and
	// verify against the header
	result := BlakeHash(header.GetHashNoNonce().Bytes(), header.GetNonce().Uint64())

	target := new(big.Int).Div(maxUint256, header.GetDifficulty())
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}

	return nil
}

// Prepare initializes the difficulty and
// total difficulty fields of a header to
// conform to the protocol
func (b *Blakimoto) Prepare(chain types.ChainReader, header types.Header) error {

	// Get the header of the block's parent.
	parent, err := chain.GetHeaderByHash(header.GetParentHash())
	if err != nil {
		if err != core.ErrBlockExists {
			return err
		}
		return ErrUnknownParent
	}

	header.SetDifficulty(b.CalcDifficulty(header, parent))
	header.SetTotalDifficulty(new(big.Int).Add(parent.GetTotalDifficulty(),
		header.GetDifficulty()))
	return nil
}
