# minecraft

This repository keeps a CLI tool to manage Minecraft servers setup to AWS EC2 instances, using Terraform.

## Why?

Why not? I wanted to learn Terraform and AWS, and what better way to do it than with Minecraft?

## Getting Started

### 1. Init terraform

```bash
make init
```

### 2. Build & Install CLI

```bash
make install
```

### 3. Run CLI

```bash
mc --help
```

## How it works?

### Modpacks

Modpacks should be written to `./modpacks` directory, with the following structure:

```markdown
modpacks/
├── my-mockpack/
│   ├── mods/
│   │   ├── mod1.jar
│   │   ├── mod2.jar
│   │   └── ...
│   └── config/
│       ├── mod1.cfg
│       ├── mod2.cfg
│       └── ...
```

### Backups

Backups are written to `./backups`.

## Forge curated modpacks examples

- [DeceasedCraft](https://www.curseforge.com/minecraft/modpacks/deceasedcraft/download/7623174)
