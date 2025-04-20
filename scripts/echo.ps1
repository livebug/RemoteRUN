
# 输入参数分钟， 持续输出时间戳，间隔1秒
# 例： .\echo.ps1 5
param(
    [int]$minutes = 1
)

echo "echo.ps1: $minutes minutes"

$endTime = (Get-Date).AddMinutes($minutes)
while ((Get-Date) -lt $endTime) {
    Write-Host (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
    Start-Sleep -Seconds 1
}
# 例： .\echo.ps1 5
