# 🛠 Go Development Standards & Tooling

This document outlines the tools and best practices we use for Go development to ensure consistency, quality, and maintainability across the team.

---

## 🧹 Code Formatting

**Tool:** `gofmt` (Built-in)  
Ensures consistent code formatting across all editors and contributors.  
✅ Always run `gofmt` before committing code.

---

## 🔍 Linting

**Tool:** `golint` CLI  
Detects:

- Bad coding practices
- Style issues
- Potential bugs

---

## 📦 Dependency Management

**Tool:** `go mod` (Built-in)  
Manages module versions and dependencies.  
➡️ Run `go mod tidy` regularly to clean up unused packages.

---

## 🧪 Unit Testing

- **Tool:** `go test` (Built-in) – Standard Go testing framework
- **Testify:** Rich assertion library (`require`, `assert`, etc.)
- **Uber's Mock:** Interface mocking for unit isolation

---

## 🧠 Static Code Analysis

**Tool:** `go vet` (Built-in)  
Detects:

- Suspicious constructs
- Common bugs
- Misuse of language features

---

## 🔐 Security Analysis

**Tool:** `gosec`  
Scans for:

- Hardcoded secrets
- Unsafe code patterns
- Cryptographic weaknesses

---

## 🛡 Vulnerability Scanning

**Tool:** `govulncheck`  

- Checks all dependencies against known CVEs
- Uses Go’s official vulnerability database  
✅ Helps keep your dependencies secure

---

## 📚 Documentation

**Tool:** `godoc`  

- Runs a local server for browsing code documentation
- Essential for onboarding and team collaboration

---

## 🔧 Environment Variable Management

> ⚠️ Avoid using `.env` files for storing sensitive information, especially in production.

### Best Practices by Environment

- **Docker:** Use `ENV` directives or docker-compose environment variables
- **AWS:** Use Secrets Manager or SSM Parameter Store
- **Kubernetes:** Use `ConfigMaps` and `Secrets`
- **Carlos/env:** Lightweight local alternative (less widely used)

---

## 🧱 Local Development Environment

- **Docker:**
  - Separate configs for **local** and **production**
  - Production and staging share the same setup

- **Air:**
  - Enables **hot reloading**
  - Improves local development feedback loop

---

By following these tools and practices, we ensure secure, scalable, and efficient Go development across the entire team.
