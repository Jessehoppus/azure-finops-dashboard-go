# 🌤️ Azure FinOps Dashboard (Go)

CLI para análise **FinOps** em **Microsoft Azure** usando **Cost Management Query API** e inventário de **Compute**.

Inspirado em [aws-finops-dashboard-go](https://github.com/Jessehoppus/aws-finops-dashboard-go).

---

## ⚙️ Requisitos

- Go 1.22+
- App Registration / Service Principal com **Client ID / Secret** ou `az login`
- **RBAC** na subscription:
  - Para custos: **Cost Management Reader**
  - Para auditorias (Compute): **Reader** na subscription
- Variáveis de ambiente (se usar SP):
  ```powershell
  $env:AZURE_TENANT_ID       = "<TENANT_ID>"
  $env:AZURE_CLIENT_ID       = "<CLIENT_ID>"
  $env:AZURE_CLIENT_SECRET   = "<CLIENT_SECRET>"
  $env:AZURE_SUBSCRIPTION_ID = "<SUBSCRIPTION_ID>"
  ```

> Dica: carregue com `.\env.ps1` (incluso no repo).

---

## 🚀 Build rápido

```powershell
.\build.ps1
# binário em .\bin\azure-finops.exe
```

Build manual:
```powershell
go mod tidy
go build -o .\bin\azure-finops.exe .\cmd\azure-finops
```

---

## 🧭 Uso

O binário aceita um **modo** inicial (`trend`, `details`, `audit`) ou pode ser chamado direto (equivalente a `trend`).

### 1) Tendência (trend)
Agrupa custos por dimensão com granularidade opcional (None/Daily/Monthly).

```powershell
.\bin\azure-finops.exe trend `
  --scope "/subscriptions/$($env:AZURE_SUBSCRIPTION_ID)" `
  --from 2025-09-01 --to 2025-10-24 `
  --dimension ServiceName `
  --granularity Monthly
```

Exportar CSV/JSON + gerar gráfico HTML:
```powershell
.\bin\azure-finops.exe trend `
  --scope "/subscriptions/$($env:AZURE_SUBSCRIPTION_ID)" `
  --from 2025-09-01 --to 2025-10-24 `
  --dimension ServiceName `
  --granularity Monthly `
  --export csv `
  --out .\reports `
  --chart .\reports\trend.html
```

Dimensões comuns: `ServiceName`, `MeterCategory`, `ResourceGroup`, `TagKey:<chave>` (ex.: `TagKey:Environment`).

### 2) Detalhes (amortizado vs. real)
Resumo por `PricingModel` + `ChargeType` (granularidade Monthly).

```powershell
.\bin\azure-finops.exe details `
  --scope "/subscriptions/$($env:AZURE_SUBSCRIPTION_ID)" `
  --from 2025-09-01 --to 2025-10-24 `
  --export json `
  --out .\reports
```

### 3) Auditorias (Compute)
- **Discos órfãos** (managed disks sem VM associada)
- **VMs paradas/deallocated** com **discos Premium**

```powershell
.\bin\azure-finops.exe audit `
  --scope "/subscriptions/$($env:AZURE_SUBSCRIPTION_ID)" `
  --export csv `
  --out .\reports
```

Saídas:
- `.\reports\audit_orphaned_disks.csv|json`
- `.\reports\audit_stopped_vms_premium.csv|json`

> Dê ao SP/usuário **Reader** na subscription para listar Compute.

---

## 🧱 Estrutura

```
cmd/azure-finops/       → CLI (modes: trend, details, audit)
internal/adapters/azure/
  ├── costquery/        → Cost Management Query API
  └── inventory/        → Compute (VMs, Disks) p/ auditorias
internal/report/        → CSV, JSON e HTML (gráfico)
env.ps1                 → carrega AZURE_* na sessão
build.ps1               → build simplificado
```

---

## 🔐 App Registration (opcional)
Crie com Azure CLI (ou use `az login` para testar via **AzureCLICredential**). Dê ao app **Cost Management Reader** na subscription e **Reader** para auditorias de Compute.

> Se quiser automatizar, posso incluir um script que crie o App, gere o secret e aplique RBAC.

---

## 🧪 CI (GitHub Actions)
- `ci.yml`: build e testes a cada push/PR
- `release.yml`: cria **release** com binários para Windows e Linux

Para publicar um release:
1. Faça um tag `git tag v0.2.0 && git push origin v0.2.0`
2. O workflow gera artefatos e anexa no release automaticamente.

---

## 📜 Licença
MIT (ou a que você preferir).
