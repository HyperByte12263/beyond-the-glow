# Downloads Malaysia's trade in HS 3304 (beauty, make-up and skincare preparations)
# from the UN Comtrade public API, one year per request (the free API is rate-limited).
#
#   powershell -File scripts/fetch_comtrade.ps1 -Flow M -From 2015 -To 2024   # imports
#   powershell -File scripts/fetch_comtrade.ps1 -Flow X -From 2024 -To 2024   # exports
#
# Output: ../raw-data/comtrade/<imports|exports>_3304_<From>_<To>.csv  (refYear, partnerCode, primaryValue in US$)
param(
    [ValidateSet('M', 'X')][string]$Flow = 'M',
    [int]$From = 2015,
    [int]$To = 2024
)

$name = if ($Flow -eq 'M') { 'imports' } else { 'exports' }
$outDir = Join-Path $PSScriptRoot '..\..\raw-data\comtrade'
$out = Join-Path $outDir "${name}_3304_${From}_${To}.csv"
$rows = @()

foreach ($yr in $From..$To) {
    $url = "https://comtradeapi.un.org/public/v1/preview/C/A/HS?reporterCode=458&period=$yr&cmdCode=3304&flowCode=$Flow"
    for ($try = 1; $try -le 4; $try++) {
        Start-Sleep -Seconds 15
        try {
            $r = Invoke-RestMethod -Uri $url -TimeoutSec 90
            # Keep one total per partner: all customs procedures, all transport modes, no second partner.
            $keep = $r.data | Where-Object { $_.partnerCode -ne 0 -and $_.customsCode -eq 'C00' -and $_.motCode -eq 0 -and $_.partner2Code -eq 0 }
            $rows += $keep | Select-Object refYear, partnerCode, primaryValue
            "$yr ok ($($keep.Count) partners)"
            break
        } catch { "$yr try $try failed: $($_.Exception.Message)" }
    }
}

$rows | Export-Csv $out -NoTypeInformation
"saved $($rows.Count) rows to $out"
