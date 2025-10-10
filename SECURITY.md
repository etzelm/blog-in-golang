# Security Policy

## Supported Versions

We actively support and provide security updates for the following versions:

| Version | Supported | End of Life |
| ------- | --------- | ----------- |
| >= 3.0  | ✅ | N/A |
| 2.2.x   | ✅ | 2026-01-01 |
| < 2.2   | ❌ | 2025-01-01 |

## Security Features

This application implements several security measures:

### Authentication & Authorization

- **bcrypt Password Hashing**: All passwords are hashed using bcrypt with salt
- **Google OAuth2 Integration**: Secure third-party authentication
- **Session Management**: Secure session handling and timeout
- **HTTPS Enforcement**: Automatic HTTPS via CertMagic in production

### Input Validation & Sanitization

- **Request Path Filtering**: Middleware blocks common malicious request paths
- **Input Validation**: Server-side validation for all user inputs
- **XSS Protection**: Template escaping for dynamic content
- **CSRF Protection**: Cross-site request forgery protection

### Infrastructure Security

- **Docker Security**: Non-root user in containers
- **AWS IAM**: Principle of least privilege for AWS resources
- **Environment Variables**: Sensitive data stored in environment variables
- **Dependency Management**: Regular security updates via Dependabot

## Reporting a Vulnerability

**⚠️ Please DO NOT report security vulnerabilities through public GitHub issues.**

### How to Report

1. **Email**: Send details to [etzelm@live.com](mailto:etzelm@live.com)
2. **Subject Line**: Use "SECURITY: [Brief Description]"
3. **PGP**: Available upon request for sensitive communications

### What to Include

Please include as much information as possible:

- **Type of vulnerability** (e.g., XSS, SQL injection, authentication bypass)
- **Steps to reproduce** the vulnerability
- **Potential impact** and attack scenarios
- **Affected versions** or components
- **Any proof-of-concept** code or screenshots
- **Suggested remediation** if you have ideas

### Response Timeline

We are committed to addressing security issues promptly:

| Timeline | Action |
|----------|--------|
| **24 hours** | Initial acknowledgment of report |
| **72 hours** | Initial assessment and severity classification |
| **7 days** | Detailed response with timeline for fix |
| **30 days** | Security patch released (for high/critical issues) |

### Severity Levels

| Severity | Examples | Response Time |
|----------|----------|---------------|
| **Critical** | RCE, Authentication bypass, Data breach | 1-7 days |
| **High** | XSS, CSRF, Privilege escalation | 7-14 days |
| **Medium** | Information disclosure, DoS | 14-30 days |
| **Low** | Minor information leaks | 30-60 days |

## Security Best Practices for Contributors

### Development

- **Never commit secrets** to the repository
- **Use `.env` files** for local development (git ignored)
- **Validate all inputs** on both client and server side
- **Follow OWASP guidelines** for web application security
- **Use HTTPS** for all external communications

### Dependencies

- **Keep dependencies updated** - Dependabot helps with this
- **Review security advisories** for dependencies
- **Use `go mod audit`** and `yarn audit` regularly
- **Pin dependency versions** in production

### AWS Configuration

- **Use IAM roles** instead of long-term access keys when possible
- **Enable CloudTrail** for audit logging
- **Configure S3 bucket policies** to restrict access
- **Enable encryption** for DynamoDB and S3

## Known Security Considerations

### Current Mitigations

- **Rate Limiting**: Implemented via Gin middleware
- **Input Sanitization**: HTML escaping in templates
- **SQL Injection**: Using AWS SDK with parameterized queries
- **Directory Traversal**: Static file serving restrictions

### Ongoing Improvements

- **Content Security Policy**: Planned for future release
- **Additional Rate Limiting**: Per-endpoint rate limiting
- **Audit Logging**: Comprehensive security event logging
- **Penetration Testing**: Regular security assessments

## Security Contact

- **Primary Contact**: [etzelm@live.com](mailto:etzelm@live.com)
- **Response Time**: Typically within 24 hours
- **Timezone**: US Eastern Time (EST/EDT)

## Acknowledgments

We appreciate security researchers who responsibly disclose vulnerabilities. Contributors who report valid security issues may be acknowledged in our security advisories (with permission).

---

For general support or non-security issues, please use:

- Contact form at [mitchelletzel.com/contact](https://mitchelletzel.com/contact)
- GitHub Issues for bugs and feature requests
