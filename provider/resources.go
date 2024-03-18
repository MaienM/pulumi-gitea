// Copyright 2016-2018, Pulumi Corporation.
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

package provider

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ettle/strcase"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"code.gitea.io/terraform-provider-gitea/gitea"
	"github.com/MaienM/pulumi-gitea/provider/pkg/version"
)

//go:embed cmd/pulumi-resource-gitea/bridge-metadata.json
var bridgeMetadata []byte

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	mainMod = "index" // the gitea module
)

var module_overrides = map[string]string{
	// Avoid module name 'public' which is a reserved keyword in (at least) nodejs.
	"Public": "PublicKey",
	// Avoid naming clash between namespace Repository and class of the same name in java.
	"Repository": "Repositories",
}

var name_overrides = map[string]string{
	// Avoid same name for type and contained field for dotnet.
	"Key": "DeploymentKey",
	"Token": "AccessToken",
}

func convertName(tfname string) (module string, name string) {
	tfNameItems := strings.Split(tfname, "_")
	contract.Assertf(len(tfNameItems) >= 2, "Invalid snake case name %s", tfname)
	contract.Assertf(tfNameItems[0] == "gitea", "Invalid snake case name %s. Does not start with gitea", tfname)
	if len(tfNameItems) == 2 {
		module = mainMod
		name = tfNameItems[1]
	} else {
		module = strcase.ToPascal(strings.Join(tfNameItems[1:len(tfNameItems)-1], "_"))
		name = tfNameItems[len(tfNameItems)-1]

		if v, ok := module_overrides[module]; ok {
			module = v
		}
	}
	contract.Assertf(!unicode.IsDigit(rune(module[0])), "Pulumi namespace must not start with a digit: %s", name)
	name = strcase.ToPascal(name)
	if v, ok := name_overrides[name]; ok {
		name = v
	}
	contract.Assertf(!unicode.IsDigit(rune(name[0])), "Pulumi name must not start with a digit: %s", name)
	return
}

func makeDataSource(ds string) tokens.ModuleMember {
	mod, name := convertName(ds)
	return tfbridge.MakeDataSource("gitea", mod, "get"+name)
}

func makeResource(res string) tokens.Type {
	mod, name := convertName(res)
	return tfbridge.MakeResource("gitea", mod, name)
}

func moduleComputeStrategy() tfbridge.Strategy {
	return tfbridge.Strategy{
		Resource: func(tfToken string, elem *tfbridge.ResourceInfo) error {
			elem.Tok = makeResource(tfToken)
			return nil
		},
		DataSource: func(tfToken string, elem *tfbridge.DataSourceInfo) error {
			elem.Tok = makeDataSource(tfToken)
			return nil
		},
	}
}

// preConfigureCallback is called before the providerConfigure function of the underlying provider.
// It should validate that the provider can be configured, and provide actionable errors in the case
// it cannot be. Configuration variables can be read from `vars` using the `stringValue` function -
// for example `stringValue(vars, "accessKey")`.
func preConfigureCallback(vars resource.PropertyMap, c shim.ResourceConfig) error {
	return nil
}

// Provider returns additional overlaid schema and metadata associated with the provider..
func Provider() tfbridge.ProviderInfo {
	// Instantiate the Terraform provider
	p := shimv2.NewProvider(gitea.Provider())

	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		P:    p,
		Name: "gitea",
		// DisplayName is a way to be able to change the casing of the provider
		// name when being displayed on the Pulumi registry
		DisplayName: "Gitea",
		// The default publisher for all packages is Pulumi.
		// Change this to your personal name (or a company name) that you
		// would like to be shown in the Pulumi Registry if this package is published
		// there.
		Publisher: "MaienM",
		// LogoURL is optional but useful to help identify your package in the Pulumi Registry
		// if this package is published there.
		//
		// You may host a logo on a domain you control or add an SVG logo for your package
		// in your repository and use the raw content URL for that file as your logo URL.
		LogoURL: "https://raw.githubusercontent.com/MaienM/pulumi-gitea/main/docs/gitea.png",
		// PluginDownloadURL is an optional URL used to download the Provider
		// for use in Pulumi programs
		// e.g https://github.com/org/pulumi-provider-name/releases/
		PluginDownloadURL: "github://api.github.com/MaienM/pulumi-gitea",
		Description:       "A Pulumi package for creating and managing Gitea resources",
		// category/cloud tag helps with categorizing the package in the Pulumi Registry.
		// For all available categories, see `Keywords` in
		// https://www.pulumi.com/docs/guides/pulumi-packages/schema/#package.
		Keywords:   []string{
			"pulumi",
			"gitea",
			"category/infrastructure",
		},
		License:    "Apache-2.0",
		Homepage:   "https://github.com/MaienM/pulumi-gitea",
		Repository: "https://github.com/MaienM/pulumi-gitea",
		// The GitHub Org for the provider - defaults to `terraform-providers`. Note that this
		// should match the TF provider module's require directive, not any replace directives.
		Version:   version.Version,
		GitHubOrg: "gitea",
		MetadataInfo: tfbridge.NewProviderMetadata(bridgeMetadata),
		TFProviderVersion: "0.3.0",
		UpstreamRepoPath: "./upstream",
		Config:    map[string]*tfbridge.SchemaInfo{
			// Add any required configuration here, or remove the example below if
			// no additional points are required.
			// "region": {
			// 	Type: tfbridge.MakeType("region", "Region"),
			// 	Default: &tfbridge.DefaultInfo{
			// 		EnvVars: []string{"AWS_REGION", "AWS_DEFAULT_REGION"},
			// 	},
			// },
		},
		PreConfigureCallback: preConfigureCallback,
		Resources:            map[string]*tfbridge.ResourceInfo{
			// Map each resource in the Terraform provider to a Pulumi type.
			//
			// "aws_iam_role": {
			//   Tok: makeResource(mainMod, "aws_iam_role"),
		  // },
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			// Map each data source in the Terraform provider to a Pulumi function.
			//
			// "aws_ami": {
			//	Tok: makeDataSource(mainMod, "aws_ami"),
			// },
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName: "@maienm/pulumi-gitea",

			// List any npm dependencies and their versions
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^3.0.0",
			},
			DevDependencies: map[string]string{
				"@types/node": "^10.0.0", // so we can access strongly typed node definitions.
				"@types/mime": "^2.0.0",
			},
			// See the documentation for tfbridge.OverlayInfo for how to lay out this
			// section, or refer to the AWS provider. Delete this section if there are
			// no overlay files.
			//Overlay: &tfbridge.OverlayInfo{},
		},
		Python: &tfbridge.PythonInfo{
			PackageName: "pulumi_gitea",

			// List any Python dependencies and their version ranges
			Requires: map[string]string{
				"pulumi": ">=3.0.0,<4.0.0",
			},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/MaienM/pulumi-%[1]s/sdk/", "gitea"),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				"gitea",
			),
			GenerateResourceContainerTypes: true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RootNamespace: "MaienM",

			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
		Java: &tfbridge.JavaInfo{
			BasePackage: "com.maienm",
		},
	}

	prov.MustComputeTokens(moduleComputeStrategy())
	prov.SetAutonaming(255, "-")

	return prov
}
