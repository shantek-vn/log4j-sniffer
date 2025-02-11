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

package java

import (
	"archive/zip"
	"bytes"
	md52 "crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"math"
	"strings"
)

type ClassHash struct {
	ClassSize               int64
	CompleteHash            string
	BytecodeInstructionHash string
}

func HashClass(jarFile string, className string) (ClassHash, error) {
	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return ClassHash{}, err
	}

	classLocation := strings.ReplaceAll(className, ".", "/")

	completeHash, size, err := md5Class(r, classLocation)
	if err != nil {
		return ClassHash{}, err
	}

	bytecodeHash, err := md5Bytecode(r, classLocation)
	if err != nil {
		return ClassHash{}, err
	}

	return ClassHash{
		ClassSize:               size,
		CompleteHash:            completeHash,
		BytecodeInstructionHash: bytecodeHash,
	}, nil
}

func ReadMethodByteCode(jarFile string, className string) (bytecode [][]byte, err error) {
	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cErr := r.Close(); err == nil && cErr != nil {
			err = cErr
		}
	}()

	classLocation := strings.ReplaceAll(className, ".", "/")

	buf, err := readClassBytes(r, classLocation)
	if err != nil {
		return nil, err
	}

	return ExtractBytecode(buf.Bytes())
}

type AverageJavaNameSizes struct {
	PackageName float32
	ClassName   float32
}

// AveragePackageAndClassLength produces an average of the package and class name lengths.
// In the event that the given Jar file contents are incredibly large, and either the total
// number of classes or packages, or the sums of their lengths, exceeds a uint32 then the
// result will be inaccurate.
//
// In the case of class or package count exceeding the maximum the average will be below the real value.
// In the case of class or package name length exceeding the maximum the average will be above the real value.
// If both values exceed the maximum then it is not possible to say whether the average will be above or below
// the real value.
//
// In practice we can handle any realistic Jar without hitting these limits.
func AveragePackageAndClassLength(files []*zip.File) AverageJavaNameSizes {
	packageNames := make(map[string]struct{})
	var classesFound, totalClassesNameSize uint32 = 0, 0
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".class") {
			parts := strings.Split(file.Name, "/") // Zip/jar is always /
			for _, packageName := range parts[:len(parts)-1] {
				packageNames[packageName] = struct{}{}
			}
			className := parts[len(parts)-1]
			classesFound = addAvoidingOverflow(classesFound, 1)
			totalClassesNameSize = addAvoidingOverflow(totalClassesNameSize, len(className)-5) // Don't count .class
		}
	}

	var packagesFound, totalPackagesNameSize uint32 = 0, 0
	for uniquePackageName := range packageNames {
		packagesFound = addAvoidingOverflow(packagesFound, 1)
		totalPackagesNameSize = addAvoidingOverflow(totalPackagesNameSize, len(uniquePackageName))
	}

	return AverageJavaNameSizes{
		PackageName: average(totalPackagesNameSize, packagesFound),
		ClassName:   average(totalClassesNameSize, classesFound),
	}
}

func average(totalSize, numberFound uint32) float32 {
	if numberFound == 0 {
		return 0
	}
	average := float32(totalSize) / float32(numberFound)
	return average
}

func addAvoidingOverflow(left uint32, right int) uint32 {
	if left > math.MaxInt32-uint32(right) {
		return math.MaxInt32
	}
	return left + uint32(right)
}

func md5Class(r *zip.ReadCloser, classLocation string) (string, int64, error) {
	c, err := r.Open(classLocation + ".class")
	if err != nil {
		return "", 0, err
	}

	h, size, err := md5File(c)
	if err != nil {
		return "", 0, err
	}

	if err := c.Close(); err != nil {
		return "", 0, err
	}
	return h, size, nil
}

func md5File(file fs.File) (string, int64, error) {
	h := md52.New()
	size, err := io.Copy(h, file)
	if err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), size, nil
}

func md5Bytecode(r *zip.ReadCloser, classLocation string) (string, error) {
	buf, err := readClassBytes(r, classLocation)
	if err != nil {
		return "", err
	}

	return HashClassInstructions(buf.Bytes())
}

func readClassBytes(r *zip.ReadCloser, classLocation string) (*bytes.Buffer, error) {
	c, err := r.Open(classLocation + ".class")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(c); err != nil {
		return nil, err
	}

	if err = c.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}
