---
title: Gitea Installation & Configuration
meta_desc: Information on how to install the Gitea provider.
layout: installation
---

## Installation

The Pulumi Gitea provider is available as a package in all Pulumi languages:

* JavaScript/TypeScript: [`@maienm/pulumi-gitea`](https://www.npmjs.com/package/@maienm/pulumi-gitea)
* Python: [`pulumi_gitea`](https://pypi.org/project/pulumi_gitea/)
* Go: [`github.com/MaienM/pulumi-gitea/sdk/go/gitea`](https://pkg.go.dev/github.com/MaienM/pulumi-gitea/sdk/go/gitea)
* .NET: [`MaienM.Gitea`](https://www.nuget.org/packages/MaienM.Gitea)


## Configuration

> Note:  
> Replace the following **sample content**, with the configuration options
> of the wrapped Terraform provider and remove this note.

The following configuration points are available for the `gitea` provider:

- `gitea:apiKey` (environment: `gitea_API_KEY`) - the API key for `gitea`
- `gitea:region` (environment: `gitea_REGION`) - the region in which to deploy resources

### Provider Binary

The Gitea provider binary is a third party binary. It can be installed using the `pulumi plugin` command.

```bash
pulumi plugin install resource gitea <version>
```

Replace the version string `<version>` with your desired version.
