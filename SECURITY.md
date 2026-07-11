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

- Change the default admin credentials (`admin@boilerworks.dev` / `password`) immediately
- Set a strong, random `SESSION_KEY` (32+ characters) — never ship the default
- Change the default database credentials in `DATABASE_URL`
- Set `ENVIRONMENT=production`
- Serve over HTTPS and add the `Secure` attribute to session and CSRF cookies
