/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package freeport_test

import (
	"net"
	"strconv"
	"testing"

	"carbonaut.dev/pkg/util/freeport"
)

func TestGetFreePort(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Error(err)
	}
	if port == 0 {
		t.Error("port == 0")
	}

	l, err := net.Listen("tcp", "localhost"+":"+strconv.Itoa(port))
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
}
