# License Choice Guide

## Overview

Open-source licenses define how others can use, modify, distribute, and contribute to your code. This guide explains the major categories, differences between common licenses, and when to choose each one.

---

# License Categories

## 1. **Permissive Licenses**

Allow broad usage with minimal obligations. These are the most popular today.

### **MIT License**

**TL;DR:** "Do whatever you want, just keep attribution and don’t sue me."

**Features:**

* Very simple
* Allows commercial & closed-source use
* Requires attribution
* No explicit patent protection

**Best for:** small libraries, examples, general OSS.

---

### **Apache License 2.0**

**TL;DR:** MIT + patent protection + contributor rights.

**Features:**

* Commercial & closed-source friendly
* Strong patent protections
* Explicit patent grant from contributors
* NOTICE file for attribution handling
* Enterprise-preferred license

**Best for:** professional libraries, corporate OSS, infrastructure tools.

---

### **BSD 2- & 3-Clause**

**TL;DR:** Similar to MIT; academic flavor.

**Features:**

* Highly permissive
* 3-clause adds “no commercial endorsement” clause

**Best for:** research-driven or historically BSD-related projects.

---

## 2. **Copyleft Licenses**

Require derivative work to remain open-source.

### **GPL (v2/v3)**

**TL;DR:** If you distribute a modified version, it must also be open-source.

**Features:**

* Strong copyleft
* Often incompatible with closed-source projects
* GPLv3 adds anti-tivoization + patent improvements

### **AGPL**

**TL;DR:** Closes the SaaS loophole — cloud-hosted modifications must be open-source.

**Best for:** tools meant to stay free even in web-service form.

### **LGPL**

**TL;DR:** Copyleft for libraries; allows linking from closed-source code.

---

## 3. **Public Domain & Public-Domain-Like**

### **Unlicense / CC0**

**TL;DR:** Minimal restrictions; essentially public domain.

**Features:**

* No attribution required
* Few legal protections

**Best for:** short examples, educational snippets.

---

# Apache 2.0 vs MIT — Quick Comparison

| Feature                       | MIT            | Apache 2.0        |
| ----------------------------- | -------------- | ----------------- |
| Commercial use                | ✔️             | ✔️                |
| Closed-source compatibility   | ✔️             | ✔️                |
| Attribution required          | ✔️             | ✔️                |
| Patent license                | ❌              | ✔️                |
| Contributor patent protection | ❌              | ✔️                |
| NOTICE file                   | ❌              | ✔️                |
| Simplicity                    | ✔️ very simple | Medium complexity |

---

# Recommendations

### ✔️ If you want **maximum adoption**:

**MIT** or **BSD 2-Clause**

### ✔️ If you want **enterprise-friendly** + **patent protection**:

**Apache License 2.0**

### ✔️ If you want the project to **stay open always**:

**GPL v3**

### ✔️ If you want to protect against closed-source SaaS forks:

**AGPL**

### ✔️ If the repo is for examples, POCs, experiments, or learning:

**MIT** (simple) or **Apache 2.0** (safer for patents)

---

# Summary

* **MIT** → simple, permissive, widely used.
* **Apache 2.0** → permissive + patent protection + structured attribution.
* **GPL/AGPL** → ensures derivative work remains open.
* **Unlicense/CC0** → public-domain style.

Choose based on your goals: adoption, protection, or openness enforcement.
