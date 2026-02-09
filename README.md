```
  ____ ____ _____ _   _ ___
/ ____|  _ \_   _| | | |_ _|
| |   | |_) || | | | | || |
| |___|  _ < | | | |_| || |
\_____|_| \_\|_|  \___/|___|
```

MANIFESTS

`application/vnd.oci.image.index.v1+json`

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 1417,
      "digest": "sha256:841c244804eb8915bea407cdefb4a5875790698b1e6cad7ad778b93e9bf60230",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 1417,
      "digest": "sha256:f3355517c51da31fb1b5c808bfd85d605959403dad01ee8dfafc16e0825ff6c8",
      "platform": {
        "architecture": "arm64",
        "os": "linux"
      }
    }
  ],
  "annotations": {
    "org.opencontainers.image.base.digest": "sha256:ae56266033b99b1ba7b360db142afc353213ff1f422b8f5f372d6207a4ec3fcf",
    "org.opencontainers.image.base.name": "cgr.dev/chainguard/static:latest"
  }
}
```

`application/vnd.oci.image.manifest.v1+json`

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:aebcefb937b8ffa0ac06820b2a52c3637aeccac32c9d78dda433f5e2d4f79b6e"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 636584,
      "digest": "sha256:60c5d5ed447870bc716279b442bb4140579b2ebc7115b185da87627f911afd10"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 127,
      "digest": "sha256:250c06f7c38e52dc77e5c7586c3e40280dc7ff9bb9007c396e06d96736cf8542"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 11019552,
      "digest": "sha256:b27e3a12f6f9712b14f05586a44e96db0d1ff9f5f45d28ffddee0eef283528aa"
    }
  ],
  "annotations": {
    "dev.chainguard.image.title": "static",
    "dev.chainguard.package.main": "",
    "org.opencontainers.image.authors": "Chainguard Team https://www.chainguard.dev/",
    "org.opencontainers.image.base.digest": "sha256:dae14bd22ea8a91c94e1bc33bdd82cc49a96cad0e711e6e03f9ae67f28a1b0be",
    "org.opencontainers.image.base.name": "cgr.dev/chainguard/static:latest",
    "org.opencontainers.image.created": "2026-01-30T05:36:52Z",
    "org.opencontainers.image.source": "https://github.com/chainguard-images/images/tree/main/images/static",
    "org.opencontainers.image.title": "static",
    "org.opencontainers.image.url": "https://images.chainguard.dev/directory/image/static/overview",
    "org.opencontainers.image.vendor": "Chainguard"
  }
}
```

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:43e2cb2da435335460e4bfe72268efca3892b48170ca4eeee4a7b229c615bdd5"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 636584,
      "digest": "sha256:60c5d5ed447870bc716279b442bb4140579b2ebc7115b185da87627f911afd10"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 127,
      "digest": "sha256:250c06f7c38e52dc77e5c7586c3e40280dc7ff9bb9007c396e06d96736cf8542"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 11019552,
      "digest": "sha256:9b26a3dd05d8db478b5c380ca606ced59ee9ba6558d81c0f2d1d4817e9082a05"
    }
  ],
  "annotations": {
    "dev.chainguard.image.title": "static",
    "dev.chainguard.package.main": "",
    "org.opencontainers.image.authors": "Chainguard Team https://www.chainguard.dev/",
    "org.opencontainers.image.base.digest": "sha256:dae14bd22ea8a91c94e1bc33bdd82cc49a96cad0e711e6e03f9ae67f28a1b0be",
    "org.opencontainers.image.base.name": "cgr.dev/chainguard/static:latest",
    "org.opencontainers.image.created": "2026-01-30T05:36:52Z",
    "org.opencontainers.image.source": "https://github.com/chainguard-images/images/tree/main/images/static",
    "org.opencontainers.image.title": "static",
    "org.opencontainers.image.url": "https://images.chainguard.dev/directory/image/static/overview",
    "org.opencontainers.image.vendor": "Chainguard"
  }
}
```
