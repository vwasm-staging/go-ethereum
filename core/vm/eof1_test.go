// Copyright 2021 The go-ethereum Authors
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

package vm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type eof1Test struct {
	code     string
	codeSize uint16
	dataSize uint16
}

var eof1ValidHeaderTests = []eof1Test{
	{"EFCAFE01010002006000", 2, 0},
	{"EFCAFE01010002020004006000AABBCCDD", 2, 4},
}

type eof1InvalidTest struct {
	code  string
	error string
}

var eof1InvalidTests = []eof1InvalidTest{
	// {"EFCAFE01010002020004006000AABBCCDD", nil}, // valid
	{"", ErrEOF1InvalidFormatByte.Error()},
	{"FE", ErrEOF1InvalidFormatByte.Error()},                                                  // invalid first byte
	{"FECAFE01010002020004006000AABBCCDD", ErrEOF1InvalidFormatByte.Error()},                  // valid except first byte
	{"EF", ErrEOF1InvalidMagic.Error()},                                                       // no magic
	{"EFCA", ErrEOF1InvalidMagic.Error()},                                                     // not complete magic
	{"EFCAFF", ErrEOF1InvalidMagic.Error()},                                                   // not correct magic
	{"EFCAFF01010002020004006000AABBCCDD", ErrEOF1InvalidMagic.Error()},                       // valid except magic
	{"EFCAFE", ErrEOF1InvalidVersion.Error()},                                                 // no version
	{"EFCAFE00", ErrEOF1InvalidVersion.Error()},                                               // invalid version
	{"EFCAFE02", ErrEOF1InvalidVersion.Error()},                                               // invalid version
	{"EFCAFE00010002020004006000AABBCCDD", ErrEOF1InvalidVersion.Error()},                     // valid except version
	{"EFCAFE01", ErrEOF1CodeSectionMissing.Error()},                                           // no header
	{"EFCAFE0100", ErrEOF1CodeSectionMissing.Error()},                                         // no code section
	{"EFCAFE0101", ErrEOF1CodeSectionSizeMissing.Error()},                                     // no code section size
	{"EFCAFE010100", ErrEOF1CodeSectionSizeMissing.Error()},                                   // code section size incomplete
	{"EFCAFE01010002", ErrEOF1InvalidTotalSize.Error()},                                       // no section terminator
	{"EFCAFE0101000200", ErrEOF1InvalidTotalSize.Error()},                                     // no code section contents
	{"EFCAFE010100020060", ErrEOF1InvalidTotalSize.Error()},                                   // not complete code section contents
	{"EFCAFE01010002006000DEADBEEF", ErrEOF1InvalidTotalSize.Error()},                         // trailing bytes after code
	{"EFCAFE010100020100020060006000", ErrEOF1MultipleCodeSections.Error()},                   // two code sections
	{"EFCAFE0101000000", ErrEOF1EmptyCodeSection.Error()},                                     // 0 size code section
	{"EFCAFE0101000002000200AABB", ErrEOF1EmptyCodeSection.Error()},                           // 0 size code section, with non-0 data section
	{"EFCAFE0102000401000200AABBCCDD6000", ErrEOF1DataSectionBeforeCodeSection.Error()},       // data section before code section
	{"EFCAFE01020004AABBCCDD", ErrEOF1DataSectionBeforeCodeSection.Error()},                   // data section without code section
	{"EFCAFE0101000202", ErrEOF1DataSectionSizeMissing.Error()},                               // no data section size
	{"EFCAFE010100020200", ErrEOF1DataSectionSizeMissing.Error()},                             // data section size incomplete
	{"EFCAFE01010002020004", ErrEOF1InvalidTotalSize.Error()},                                 // no section terminator
	{"EFCAFE01010002020004006000", ErrEOF1InvalidTotalSize.Error()},                           // no data section contents
	{"EFCAFE01010002020004006000AABBCC", ErrEOF1InvalidTotalSize.Error()},                     // not complete data section contents
	{"EFCAFE01010002020004006000AABBCCDDEE", ErrEOF1InvalidTotalSize.Error()},                 // trailing bytes after data
	{"EFCAFE01010002020000006000", ErrEOF1EmptyDataSection.Error()},                           // 0 size data section
	{"EFCAFE01010002020004020004006000AABBCCDDAABBCCDD", ErrEOF1MultipleDataSections.Error()}, // two data sections
	{"EFCAFE01010002030004006000AABBCCDD", ErrEOF1UnknownSection.Error()},                     // section id = 3
}

func TestReadEOF1Header(t *testing.T) {

	for _, test := range eof1ValidHeaderTests {
		header, err := readEOF1Header(common.Hex2Bytes(test.code))
		if err != nil {
			t.Errorf("code %v validation failure, error: %v", test.code, err)
		}
		if header.codeSize != test.codeSize {
			t.Errorf("code %v codeSize expected %v, got %v", test.code, test.codeSize, header.codeSize)
		}
		if header.dataSize != test.dataSize {
			t.Errorf("code %v dataSize expected %v, got %v", test.code, test.dataSize, header.dataSize)
		}
	}
}

func TestIsValidEOF(t *testing.T) {

	for _, test := range eof1InvalidTests {
		if isValidEOF(common.Hex2Bytes(test.code)) {
			t.Errorf("code %v expected to be invalid", test.code)
		}
	}
}

func TestValidateEOF1(t *testing.T) {

	for _, test := range eof1InvalidTests {
		_, err := readEOF1Header(common.Hex2Bytes(test.code))
		if err == nil {
			t.Errorf("code %v expected to be invalid", test.code)
		}
		if err.Error() != test.error {
			t.Errorf("code %v expected error: \"%v\" got error: \"%v\"", test.code, test.error, err.Error())
		}
	}
}
