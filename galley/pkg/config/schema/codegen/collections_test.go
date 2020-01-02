// Copyright 2019 Istio Authors
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

package codegen

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"istio.io/istio/galley/pkg/config/schema/ast"
)

func TestStaticCollections(t *testing.T) {
	var cases = []struct {
		packageName string
		m           *ast.Metadata
		err         string
		output      string
	}{
		{
			packageName: "pkg",
			m: &ast.Metadata{
				Collections: []*ast.Collection{
					{
						Name:         "foo",
						VariableName: "Foo",
						Description:  "describes a really cool foo thing",
						Group:        "foo.group",
						Kind:         "fookind",
						Disabled:     true,
					},
					{
						Name:         "bar",
						VariableName: "Bar",
						Description:  "describes a really cool bar thing",
						Group:        "bar.group",
						Kind:         "barkind",
						Disabled:     false,
					},
				},
				Resources: []*ast.Resource{
					{
						Group:         "foo.group",
						Version:       "v1",
						Kind:          "fookind",
						Plural:        "fookinds",
						ClusterScoped: true,
						Proto:         "google.protobuf.Struct",
						ProtoPackage:  "github.com/gogo/protobuf/types",
						Validate:      "EmptyValidate",
					},
					{
						Group:         "bar.group",
						Version:       "v1",
						Kind:          "barkind",
						Plural:        "barkinds",
						ClusterScoped: false,
						Proto:         "google.protobuf.Struct",
						ProtoPackage:  "github.com/gogo/protobuf/types",
						Validate:      "EmptyValidate",
					},
				},
			},
			output: `
// GENERATED FILE -- DO NOT EDIT
//

package pkg

import (
	"istio.io/istio/galley/pkg/config/schema/collection"
	"istio.io/istio/galley/pkg/config/schema/resource"
	"istio.io/istio/pkg/config/validation"
)

var (

	// Bar describes a really cool bar thing
	Bar = collection.Builder {
		Name: "bar",
		Disabled: false,
		Schema: resource.Builder {
			Group: "bar.group",
			Kind: "barkind",
			Plural: "barkinds",
			Version: "v1",
			Proto: "google.protobuf.Struct",
			ProtoPackage: "github.com/gogo/protobuf/types",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// Foo describes a really cool foo thing
	Foo = collection.Builder {
		Name: "foo",
		Disabled: true,
		Schema: resource.Builder {
			Group: "foo.group",
			Kind: "fookind",
			Plural: "fookinds",
			Version: "v1",
			Proto: "google.protobuf.Struct",
			ProtoPackage: "github.com/gogo/protobuf/types",
			ClusterScoped: true,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()


	// All contains all collections in the system.
	All = collection.NewSchemasBuilder().
		MustAdd(Bar).
		MustAdd(Foo).
		Build()

	// Istio contains only Istio collections.
	Istio = collection.NewSchemasBuilder().
		Build()

	// Kube contains only kubernetes collections.
	Kube = collection.NewSchemasBuilder().
		Build()
)
`,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			g := NewGomegaWithT(t)

			s, err := StaticCollections(c.packageName, c.m)
			if c.err != "" {
				g.Expect(err).NotTo(BeNil())
				g.Expect(err.Error()).To(Equal(s))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(strings.TrimSpace(s)).To(Equal(strings.TrimSpace(c.output)))
			}
		})
	}
}
