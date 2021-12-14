// Copyright (c) 2020 Tigera, Inc. All rights reserved.
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

package ut_test

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/projectcalico/calico/felix/bpf"
	"github.com/projectcalico/calico/felix/bpf/conntrack"
)

func TestMapEntryDeletion(t *testing.T) {
	k, err1 := setUpMapTestWithSingleKV(t)

	err2 := ctMap.Delete(k.AsBytes())
	err3 := ctMap.Delete(k.AsBytes())

	// Defer error checking since the Delete calls do the cleanup for this test...
	Expect(err1).NotTo(HaveOccurred(), "Failed to create map entry")
	Expect(err2).NotTo(HaveOccurred(), "Failed to delete map entry")
	Expect(bpf.IsNotExists(err3)).To(Equal(true), "Error from deletion of non-existent entry was incorrect")
}

func TestMapIterActionDelete(t *testing.T) {
	k, err1 := setUpMapTestWithSingleKV(t)

	seenKey := false
	iterErr := ctMap.Iter(func(k2, v []byte) bpf.IteratorAction {
		seenKey = seenKey || reflect.DeepEqual(k2, k[:])
		return bpf.IterDelete
	})

	seenKeyAfterDel := false
	iterErr2 := ctMap.Iter(func(k2, v []byte) bpf.IteratorAction {
		seenKeyAfterDel = seenKeyAfterDel || reflect.DeepEqual(k2, k[:])
		return bpf.IterNone
	})

	// Defer error checking since the Delete call does the cleanup for this test...
	Expect(err1).NotTo(HaveOccurred(), "Failed to create map entry")
	Expect(iterErr).ToNot(HaveOccurred(), "Failed to iterate map")
	Expect(seenKey).To(BeTrue(), "Expected to see the key we put in the map")
	Expect(seenKeyAfterDel).To(BeFalse(), "Saw key we'd just deleted")
	Expect(iterErr2).ToNot(HaveOccurred(), "Failed to iterate map after delete")
}

func TestMapIterationDeleteAfter(t *testing.T) {
	k, err1 := setUpMapTestWithSingleKV(t)

	seenKey := false
	iterErr := ctMap.Iter(func(k2, v []byte) bpf.IteratorAction {
		seenKey = seenKey || reflect.DeepEqual(k2, k[:])
		return bpf.IterNone
	})

	err2 := ctMap.Delete(k.AsBytes())

	seenKeyAfterDel := false
	iterErr2 := ctMap.Iter(func(k2, v []byte) bpf.IteratorAction {
		seenKeyAfterDel = seenKeyAfterDel || reflect.DeepEqual(k2, k[:])
		return bpf.IterNone
	})

	// Defer error checking since the Delete call does the cleanup for this test...
	Expect(err1).NotTo(HaveOccurred(), "Failed to create map entry")
	Expect(iterErr).ToNot(HaveOccurred(), "Failed to iterate map")
	Expect(seenKey).To(BeTrue(), "Expected to see the key we put in the map")
	Expect(seenKeyAfterDel).To(BeFalse(), "Saw key we'd just deleted")
	Expect(iterErr2).ToNot(HaveOccurred(), "Failed to iterate map after delete")
	Expect(err2).NotTo(HaveOccurred(), "Failed to delete map entry")
}

func setUpMapTestWithSingleKV(t *testing.T) (conntrack.Key, error) {
	RegisterTestingT(t)
	k := conntrack.NewKey(1, net.ParseIP("10.0.0.1"), 51234, net.ParseIP("10.0.0.2"), 8080)
	v := conntrack.Value{}
	for i := range v {
		v[i] = uint8(i)
	}
	err1 := ctMap.Update(k.AsBytes(), v[:])
	return k, err1
}

func BenchmarkMapIteration10k(b *testing.B) {
	benchMapIteration(b, 10000)
}

func BenchmarkMapIteration100k(b *testing.B) {
	benchMapIteration(b, 100000)
}

func BenchmarkMapIteration500k(b *testing.B) {
	benchMapIteration(b, 500000)
}

func benchMapIteration(b *testing.B, n int) {
	defer cleanUpMaps()
	logLevel := logrus.GetLevel()
	logrus.SetLevel(logrus.InfoLevel)
	defer logrus.SetLevel(logLevel)
	setUpConntrackMapEntries(b, n)
	var keepK, keepV []byte
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := ctMap.Iter(func(k, v []byte) bpf.IteratorAction {
			keepK = k
			keepV = v
			return bpf.IterNone
		})
		if err != nil {
			panic(err)
		}
	}
	runtime.KeepAlive(keepK)
	runtime.KeepAlive(keepV)
	b.StopTimer()
}

func setUpConntrackMapEntries(b *testing.B, n int) {
	b.StopTimer()
	for i := 0; i < n; i++ {
		var k conntrack.Key
		var v conntrack.Value
		binary.LittleEndian.PutUint32(k[:], uint32(i))
		err := ctMap.Update(k.AsBytes(), v[:])
		if err != nil {
			panic(err)
		}
	}
	b.StartTimer()
}

func BenchmarkMapIteratorMulti10(b *testing.B) {
	benchMapIteratorMulti(b, 10)
}

func BenchmarkMapIteratorMulti100(b *testing.B) {
	benchMapIteratorMulti(b, 100)
}

func BenchmarkMapIteratorMulti101(b *testing.B) {
	benchMapIteratorMulti(b, 101)
}

func BenchmarkMapIteratorMulti10k(b *testing.B) {
	benchMapIteratorMulti(b, 10000)
}

func BenchmarkMapIteratorMulti100k(b *testing.B) {
	benchMapIteratorMulti(b, 100000)
}

func BenchmarkMapIteratorMulti500k(b *testing.B) {
	benchMapIteratorMulti(b, 500000)
}

func benchMapIteratorMulti(b *testing.B, n int) {
	logLevel := logrus.GetLevel()
	logrus.SetLevel(logrus.InfoLevel)
	defer logrus.SetLevel(logLevel)
	setUpConntrackMapEntries(b, n)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		doSingleMapIteratorMultiTest(n)
	}
	b.StopTimer()
	cleanUpMaps()
}

func doSingleMapIteratorMultiTest(n int) {
	iter, err := bpf.NewMapIterator(ctMap.MapFD(), conntrack.KeySize, conntrack.ValueSize, conntrack.MaxEntries)
	if err != nil {
		panic(err)
	}

	numIterations := 0
	var k, v []byte
	for {
		k, v, err = iter.Next()
		if err != nil {
			if err == bpf.ErrIterationFinished {
				break
			}
			panic(err)
		}
		numIterations++
	}
	if numIterations != n {
		panic(fmt.Sprintf("Unexpected number of iterations: %d", numIterations))
	}
	_ = iter.Close()
	runtime.KeepAlive(k)
	runtime.KeepAlive(v)
}
