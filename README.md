# 🌤️ Azure FinOps Dashboard (Go)

Ferramenta CLI para análise FinOps em ambientes **Microsoft Azure**, baseada nas APIs do **Cost Management**.

Inspirado no projeto original [aws-finops-dashboard-go](https://github.com/Jessehoppus/aws-finops-dashboard-go).

---

## ⚙️ Requisitos

- Go 1.22+
- Permissão **Cost Management Reader** na subscription/RG desejado.
- App Registration no Entra ID (Client ID, Secret, Tenant ID).
- Variáveis de ambiente configuradas (veja abaixo).

> Dica: Use o script `scripts/setup-appregistration.ps1` para criar automaticamente o App Registration, definir o papel **Cost Management Reader** e devolver os comandos de ambiente.

```powershell
# Exemplo de variáveis no PowerShell
$env:AZURE_CLIENT_ID="xxxx-..."
$env:AZURE_TENANT_ID="xxxx-..."
$env:AZURE_CLIENT_SECRET="xxxx-..."
$env:AZURE_SUBSCRIPTION_ID="xxxx-..."
```

---

## 🚀 Build e Execução

```bash
make build
./bin/azure-finops trend \
  --scope "/subscriptions/$AZURE_SUBSCRIPTION_ID" \
  --from 2025-09-24 \
  --to 2025-10-24 \
  --dimension ServiceName \
  --granularity Monthly
```

Saída (exemplo):

```
Período: 2025-09-24 .. 2025-10-24 | Escopo: /subscriptions/xxxx
ServiceName                         totalCost
Virtual Machines                     1245.87
Storage                              512.34
...
```

---

## 🧱 Estrutura

```
cmd/azure-finops/       → CLI
internal/core/           → Portas / interfaces
internal/adapters/azure/ → SDK Azure Cost Management
pkg/version/             → Metadados
scripts/                 → Automação (App Registration)
.github/workflows/       → CI (GitHub Actions)
```

---

## 🧩 Próximos passos

- Módulo `costdetails` para Savings Plans/Reservations (amortizado vs real)
- Exportação CSV/JSON
- Auditoria de recursos (VMs ociosas, discos órfãos)
- Painel gráfico (Dash/Grafana)
