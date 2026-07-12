# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Boilerworks, please report it responsibly.

**Do not open a public issue.**

Instead, email **security@weareconflict.com** with:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will acknowledge your report within 48 hours and aim to release a fix within 7 days for critical issues.

## Supported Versions

| Version | Supported |
| ------- | --------- |
| latest  | Yes       |

## Security Best Practices

When deploying Boilerworks:

- Change all default credentials (database password, `SESSION_KEY`, and the seeded `admin@boilerworks.dev` user)
- Use HTTPS in production so session cookies are protected in transit
- Set `ENVIRONMENT=production` for the Go API and `NODE_ENV=production` for the Next.js frontend
- Set `FRONTEND_URL` to your domain only — it drives the CORS allowlist
