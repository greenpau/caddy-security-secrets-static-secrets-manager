# caddy-security-secrets-static-secrets-manager

[![build](https://github.com/greenpau/caddy-security-secrets-static-secrets-manager/actions/workflows/build.yml/badge.svg)](https://github.com/greenpau/caddy-security-secrets-static-secrets-manager/actions/workflows/build.yml)
[![docs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/greenpau/caddy-security-secrets-static-secrets-manager)

[Caddy Security](https://github.com/greenpau/caddy-security) Secrets Plugin
for statically configured secrets.

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
  * [Generate Secrets](#generate-secrets)
  * [Building Caddy](#building-caddy)
  * [Caddyfile Usage](#caddyfile-usage)
    * [Without Plugin](#without-plugin)
    * [Plugin Configuration](#plugin-configuration)

<!-- end-markdown-toc -->

## Getting Started

### Generate Secrets

Please follow this [doc](https://github.com/greenpau/go-authcrunch-secrets-static-secrets-manager#getting-started)
to generate secrets.

### Building Caddy

For `secrets static_secrets_manager` directives to work, build `caddy` with the
`latest` version of this plugin.

```bash
xcaddy build ... \
  --with github.com/greenpau/caddy-security-secrets-static-secrets-manager@latest
```

### Caddyfile Usage

#### Without Plugin

The following is a snippet of `Caddyfile` without the use of this plugin.

```
{
        security {
                local identity store localdb {
                        realm local
                        path /etc/caddy/users.json
                        user jsmith {
                                name John Smith
                                email jsmith@localhost.localdomain
                                password "bcrypt:10:$2a$10$iqq53VjdCwknBSBrnyLd9OH1Mfh6kqPezMMy6h6F41iLdVDkj13I6" overwrite
                                roles authp/admin authp/user
                        }
                }

                authentication portal myportal {
                        crypto default token lifetime 3600
                        crypto key sign-verify b006d65b-c923-46a1-8da1-7d52558508fe
                        enable identity store localdb
                }
        }
}
```

#### Plugin Configuration

Now, here is the configuration using `secrets` retrieved from locally configured secrets:

```
{
	security {
		secrets static_secrets_manager access_token {
			shared_secret b006d65b-c923-46a1-8da1-7d52558508fe
		}

		secrets static_secrets_manager users/jsmith {
			name "John Smith"
			email "jsmith@localhost.localdomain"
			password "bcrypt:10:$2a$10$iqq53VjdCwknBSBrnyLd9OH1Mfh6kqPezMMy6h6F41iLdVDkj13I6"
			api_key "bcrypt:10:$2a$10$TEQ7ZG9cAdWwhQK36orCGOlokqQA55ddE0WEsl00oLZh567okdcZ6"
		}

		local identity store localdb {
			realm local
			path users.json
			user jsmith {
				name "secrets:users/jsmith:name"
				email "secrets:users/jsmith:email"
				password "secrets:users/jsmith:password" overwrite
				api_key "secrets:users/jsmith:api_key" overwrite
				roles authp/admin authp/user
			}
		}

		authentication portal myportal {
			crypto default token lifetime 3600
			crypto key sign-verify "secrets:access_token:shared_secret"
			enable identity store localdb
		}
	}
}
```
