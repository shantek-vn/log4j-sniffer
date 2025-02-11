// Copyright (c) 2021 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crawl

import (
	"bytes"
)

// Generated using log4j-sniffer identify
var classMd5s = map[string]string{
	"6b15f42c333ac39abacfeeeb18852a44": "2.1-2.3",
	"8b2260b1cce64144f6310876f94b1638": "2.4-2.5",
	"3bd9f41b89ce4fe8ccbf73e43195a5ce": "2.6-2.6.2",
	"415c13e7c8505fb056d540eac29b72fa": "2.7-2.8.1",
	"a193703904a3f18fb3c90a877eb5c8a7": "2.8.2",
	"04fdd701809d17465c17c7e603b1b202": "2.9.0-2.11.2",
	"5824711d6c68162eb535cc4dbf7485d3": "2.12.0",
	"102cac5b7726457244af1f44e54ff468": "2.12.2",
	"21f055b62c15453f0d7970a9d994cab7": "2.13.0-2.13.3",
	"f1d630c48928096a484e4b95ccb162a0": "2.14.0-2.14.1",
	"5d253e53fa993e122ff012221aa49ec3": "2.15.0",
	"ba1cf8f81e7b31c709768561ba8ab558": "2.16.0",
	"3dc5cf97546007be53b2f3d44028fa58": "2.17.0",
}

var bytecodeMd5s = map[string]string{
	"e873c1367963fad624f7128e74013725-v0": "2.1-2.5",
	"34603528cf70de0e17669acd122ad110-v0": "2.6-2.8.1",
	"bdbc07b787588e54870b5e90933d2306-v0": "2.8.2",
	"bd12d274eef8fa455f303284834ce62b-v0": "2.9.0-2.11.2",
	"81fcf4a9f7dd4dcb4fa0ab6daaed496f-v0": "2.12.2",
	"8139e14cd3955ef709139c3f23d38057-v0": "2.12.0-2.14.1",
	"5120cdf3b914bb4347e3235efce4eabf-v0": "2.15.0",
	"0761bbaeee745db2559b6416a3a30712-v0": "2.16.0",
	"79cd7e06b1a00b375f221414f06bbdd6-v0": "2.17.0",
}

type partialMethodMatchSignature struct {
	Prefix []byte
	Suffix []byte
}

type partialBytecodeSignature struct {
	Version        string
	PartialMatches []partialMethodMatchSignature
}

type exactMatch struct {
	Versions []string
	Match    []byte
}

var exactMatches = []exactMatch{
	{
		Versions: []string{"2.16.0", "2.15.0"},
		Match:    []byte{0x2a, 0x01, 0x2b, 0xb7, 0x2a, 0x2c, 0xb5, 0x2a, 0x2d, 0xb5, 0x2a, 0x19, 0xb5, 0x2a, 0x19, 0xb5, 0xb1},
	},
	{
		Versions: []string{"2.16.0"},
		Match:    []byte{0x2a, 0x01, 0x2b, 0xb7, 0x2a, 0x01, 0xb5, 0x2a, 0x01, 0xb5, 0x2a, 0x01, 0xb5, 0x2a, 0x01, 0xb5, 0xb1},
	},
	{
		Versions: []string{"2.16.0"},
		Match:    []byte{0x2a, 0xb4, 0xc6, 0x2a, 0xb4, 0xb8, 0xac, 0x04, 0xac},
	},
	{
		Versions: []string{"2.16.0", "2.15.0"},
		Match:    []byte{0x2a, 0x2b, 0x2c, 0x2d, 0x19, 0x19, 0xb7, 0xb1},
	},
	{
		Versions: []string{"2.9.0-2.14.1", "2.17.0", "2.12.2", "2.8.2"},
		Match:    []byte{0x2a, 0x01, 0x2b, 0xb7, 0x2a, 0x2c, 0xb5, 0xb1},
	},
	{
		Versions: []string{"2.9.0-2.14.1", "2.16.0", "2.15.0", "2.17.0", "2.8.2"},
		Match:    []byte{0x12, 0xb6, 0xb2, 0x01, 0xb8, 0xc0, 0xb0},
	},
	{
		Versions: []string{"2.1-2.8.1"},
		Match:    []byte{0x2a, 0x2b, 0xb7, 0x2a, 0x2c, 0xb5, 0xb1},
	},
	{
		Versions: []string{"2.9.0-2.14.1", "2.8.2", "2.1-2.8.1"},
		Match:    []byte{0x2a, 0xb4, 0x2b, 0xb9, 0xb0},
	},
	{
		Versions: []string{"2.9.0-2.14.1", "2.17.0", "2.12.2", "2.8.2", "2.1-2.8.1"},
		Match:    []byte{0x2a, 0x2b, 0x2c, 0xb7, 0xb1},
	},
	{
		Versions: []string{"2.16.0", "2.12.2"},
		Match:    []byte{0xb8, 0x12, 0x03, 0xb6, 0xac},
	},
	{
		Versions: []string{"2.9.0-2.14.1", "2.15.0", "2.17.0", "2.12.2", "2.8.2"},
		Match:    []byte{0x2a, 0xb4, 0xb8, 0xac},
	},
	{
		Versions: []string{"2.16.0"},
		Match:    []byte{0x2a, 0x2b, 0xb7, 0xb1},
	},
	{
		Versions: []string{"2.17.0"},
		Match:    []byte{0x12, 0xb8, 0xac},
	},
}

var partialBytecodeSignatures = map[string]partialBytecodeSignature{
	"2.9.0-2.14.1": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0xb6, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0xb7, 0xb3, 0xb1},
			},
		},
	},
	"2.16.0": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0x2a, 0xb4, 0xc7, 0x01, 0xb0, 0xbb, 0x59, 0x2b, 0xb7},
				Suffix: []byte{0xb2, 0x12, 0x2b, 0xb9, 0x01, 0xb0, 0x2a, 0xb4, 0x2b, 0xb9, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x53, 0x59, 0x04, 0x12, 0x53, 0x59, 0x05, 0x12, 0x53, 0xb8, 0xb3, 0xb1},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0xb6, 0xb0},
			},
		},
	},
	"2.15.0": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0xbb, 0x59, 0x2b, 0xb7},
				Suffix: []byte{0xb2, 0x12, 0x2b, 0xb9, 0x01, 0xb0, 0x2a, 0xb4, 0x2b, 0xb9, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x53, 0x59, 0x04, 0x12, 0x53, 0x59, 0x05, 0x12, 0x53, 0xb8, 0xb3, 0xb1},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0xb6, 0xb0},
			},
		},
	},
	"2.17.0": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0x2a, 0xb4, 0xc7, 0x01, 0xb0, 0xbb, 0x59, 0x2b, 0xb7},
				Suffix: []byte{0xb2, 0x12, 0x2b, 0xb9, 0x01, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0xb6, 0xb0},
			},
			{
				Prefix: []byte{0xb8, 0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb6, 0xb6, 0x03, 0xb6, 0xac},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0xb7, 0xb3, 0xb1},
			},
		},
	},
	"2.12.2": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0x2a, 0xb4, 0xb6, 0x12, 0xb6, 0xb6, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0xb7, 0xb3, 0xb1},
			},
		},
	},
	"2.8.2": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0xbb, 0x59, 0xb7, 0x12, 0xb6, 0xb6, 0x10, 0xb6, 0x12, 0xb6, 0xb6, 0xb6},
				Suffix: []byte{0xb9, 0x19, 0xc6, 0x19, 0x19, 0xb6, 0x19, 0xb2, 0x19, 0xb8, 0xc0, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0xb7, 0xb3, 0xb1},
			},
		},
	},
	"2.1-2.8.1": {
		PartialMatches: []partialMethodMatchSignature{
			{
				Prefix: []byte{0xbb, 0x59, 0xb7},
				Suffix: []byte{0xb9, 0x19, 0xc6, 0x19, 0x19, 0xb6, 0x19, 0xb2, 0x19, 0xb8, 0xc0, 0xb0},
			},
			{
				Prefix: []byte{0xbb, 0x59},
				Suffix: []byte{0xb7, 0xb3, 0xb1},
			},
			{
				Prefix: []byte{0x2a, 0xb4},
				Suffix: []byte{0xb1},
			},
		},
	},
}

// BytecodeMatchesPartialSignatures compares the given class method bytecode against snippets from known versions.
// A partial signature is made up of two parts: exact matches and partial matches.
// For an exact match to be identified the entirety of the bytecode a method must match the signature.
// Partial matches provide a prefix and suffix, these must both match a given method for the partial match to be a success.
func BytecodeMatchesPartialSignatures(methodBytecodes [][]byte) (string, bool) {
	exactVersionMatched := make(map[string]bool)
	matchedIndexes := make([]bool, len(methodBytecodes))
	for _, exactMatch := range exactMatches {
		matchIndex := -1
		for i, methodBytecode := range methodBytecodes {
			if matchedIndexes[i] {
				continue
			}
			if bytes.Compare(exactMatch.Match, methodBytecode) == 0 {
				matchIndex = i
				break
			}
		}
		for _, version := range exactMatch.Versions {
			exactVersionMatched[version] = matchIndex != -1
		}
		if matchIndex != -1 {
			matchedIndexes[matchIndex] = true
		}
	}

	for version, matched := range exactVersionMatched {
		if matched {
			partialSignature := partialBytecodeSignatures[version]
			if partialSignatureMatches(methodBytecodes, partialSignature) {
				return version, true
			}
		}
	}

	return UnknownVersion, false
}

func partialSignatureMatches(methodBytecodes [][]byte, partialSignature partialBytecodeSignature) bool {
	for _, partialMatch := range partialSignature.PartialMatches {
		matchIndex := -1
		matchedIndexes := make([]bool, len(methodBytecodes))
		for i, methodBytecode := range methodBytecodes {
			if matchedIndexes[i] {
				continue
			}
			if len(methodBytecode) < len(partialMatch.Prefix)+len(partialMatch.Suffix) {
				continue
			}
			matched := true
			for x := 0; x < len(partialMatch.Prefix); x++ {
				if partialMatch.Prefix[x] != methodBytecode[x] {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
			bytecodeLength, suffixLength := len(methodBytecode), len(partialMatch.Suffix)
			for x := 0; x < suffixLength; x++ {
				if partialMatch.Suffix[suffixLength-x-1] != methodBytecode[bytecodeLength-x-1] {
					matched = false
					break
				}
			}
			if matched {
				matchIndex = i
				break
			}
		}
		if matchIndex == -1 {
			return false
		}
		matchedIndexes[matchIndex] = true
	}
	return true
}
