<#
.SYNOPSIS
  Cria App Registration (Entra ID) + Service Principal, gera Client Secret,
  e atribui a role "Cost Management Reader" no escopo informado (subscription por padrão).

.PREREQUISITES
  - Azure CLI (az) autenticado: az login
  - Permissões para criar app/SP e conceder role assignment

.PARAMETERS
  -AppName           Nome do App Registration (default: azure-finops-dashboard-go)
  -SubscriptionId    Subscription alvo para role assignment
  -Scope             Escopo completo para role (default: /subscriptions/<SubscriptionId>)
  -SecretYears       Validade do secret em anos (default: 1)
#>

param(
  [string]$AppName = "azure-finops-dashboard-go",
  [string]$SubscriptionId,
  [string]$Scope,
  [int]$SecretYears = 1
)

if (-not $SubscriptionId) {
  $current = az account show --query id -o tsv 2>$null
  if (-not $current) {
    Write-Error "Assine com 'az login' ou informe -SubscriptionId"
    exit 1
  }
  $SubscriptionId = $current
}

if (-not $Scope) {
  $Scope = "/subscriptions/$SubscriptionId"
}

Write-Host ">>> SubscriptionId: $SubscriptionId"
Write-Host ">>> Scope: $Scope"
Write-Host ">>> AppName: $AppName"

# Tenant
$TenantId = az account show --query tenantId -o tsv
Write-Host ">>> TenantId: $TenantId"

# Create App
$app = az ad app create --display-name $AppName --query "{appId:appId, id:id}" -o json
$appObj = $app | ConvertFrom-Json
$ClientId = $appObj.appId

Write-Host ">>> App criado. ClientId: $ClientId"

# Create Service Principal
$sp = az ad sp create --id $ClientId | Out-Null
Write-Host ">>> Service Principal criado."

# Create Secret
$endDate = (Get-Date).AddYears($SecretYears).ToString("yyyy-MM-dd")
$cred = az ad app credential reset --id $ClientId --end-date $endDate --query "{clientSecret:password}" -o json
$credObj = $cred | ConvertFrom-Json
$ClientSecret = $credObj.clientSecret

# Role assignment
Write-Host ">>> Atribuindo papel 'Cost Management Reader' no escopo $Scope"
az role assignment create --assignee $ClientId --role "Cost Management Reader" --scope $Scope | Out-Null

# Output vars
Write-Host ""
Write-Host "==============================================="
Write-Host " App Registration criado com sucesso!"
Write-Host "==============================================="
Write-Host "AZURE_TENANT_ID      = $TenantId"
Write-Host "AZURE_CLIENT_ID      = $ClientId"
Write-Host "AZURE_CLIENT_SECRET  = $ClientSecret"
Write-Host "AZURE_SUBSCRIPTION_ID= $SubscriptionId"
Write-Host ""
Write-Host "PowerShell (sessão atual):"
Write-Host '$env:AZURE_TENANT_ID="'$TenantId'"'
Write-Host '$env:AZURE_CLIENT_ID="'$ClientId'"'
Write-Host '$env:AZURE_CLIENT_SECRET="'$ClientSecret'"'
Write-Host '$env:AZURE_SUBSCRIPTION_ID="'$SubscriptionId'"'
Write-Host ""
Write-Host "Bash:"
Write-Host 'export AZURE_TENANT_ID='$TenantId
Write-Host 'export AZURE_CLIENT_ID='$ClientId
Write-Host 'export AZURE_CLIENT_SECRET='$ClientSecret
Write-Host 'export AZURE_SUBSCRIPTION_ID='$SubscriptionId
Write-Host ""
Write-Host "IMPORTANTE: Guarde o Client Secret com segurança."
