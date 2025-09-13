# Bookstore API Testing Script
# This script tests all the API endpoints

$baseUrl = "http://localhost:8080"
$apiUrl = "$baseUrl/api/v1"

Write-Host "=== Bookstore API Testing Script ===" -ForegroundColor Green
Write-Host "Base URL: $baseUrl" -ForegroundColor Yellow
Write-Host "API URL: $apiUrl" -ForegroundColor Yellow
Write-Host ""

# Function to make HTTP requests
function Invoke-APIRequest {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Body = $null,
        [hashtable]$Headers = @{}
    )
    
    try {
        $params = @{
            Method = $Method
            Uri = $Url
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = "application/json"
        }
        
        if ($Headers.Count -gt 0) {
            $params.Headers = $Headers
        }
        
        $response = Invoke-WebRequest @params
        return @{
            StatusCode = $response.StatusCode
            Content = $response.Content | ConvertFrom-Json
            Success = $true
        }
    }
    catch {
        return @{
            StatusCode = $_.Exception.Response.StatusCode.value__
            Content = $_.Exception.Message
            Success = $false
        }
    }
}

# Test 1: Health Check
Write-Host "1. Testing Health Check..." -ForegroundColor Cyan
$healthResponse = Invoke-APIRequest -Method "GET" -Url "$baseUrl/health"
Write-Host "   Status: $($healthResponse.StatusCode)" -ForegroundColor White
Write-Host "   Response: $($healthResponse.Content | ConvertTo-Json -Compress)" -ForegroundColor Gray
Write-Host ""

# Test 2: API Documentation
Write-Host "2. Testing API Documentation..." -ForegroundColor Cyan
$docsResponse = Invoke-APIRequest -Method "GET" -Url "$baseUrl/docs"
Write-Host "   Status: $($docsResponse.StatusCode)" -ForegroundColor White
Write-Host "   Title: $($docsResponse.Content.data.title)" -ForegroundColor Gray
Write-Host ""

# Test 3: Get All Authors (should work without auth)
Write-Host "3. Testing Get All Authors..." -ForegroundColor Cyan
$authorsResponse = Invoke-APIRequest -Method "GET" -Url "$apiUrl/authors"
Write-Host "   Status: $($authorsResponse.StatusCode)" -ForegroundColor White
Write-Host "   Count: $($authorsResponse.Content.data.Count)" -ForegroundColor Gray
Write-Host ""

# Test 4: Create Author (should fail without auth)
Write-Host "4. Testing Create Author (without auth - should fail)..." -ForegroundColor Cyan
$authorData = @{
    name = "Test Author"
    email = "test@example.com"
    biography = "A test author"
} | ConvertTo-Json

$createAuthorResponse = Invoke-APIRequest -Method "POST" -Url "$apiUrl/authors" -Body $authorData
Write-Host "   Status: $($createAuthorResponse.StatusCode)" -ForegroundColor White
Write-Host "   Message: $($createAuthorResponse.Content.message)" -ForegroundColor Gray
Write-Host ""

# Test 5: Create Author (with fake auth - should work)
Write-Host "5. Testing Create Author (with fake auth)..." -ForegroundColor Cyan
$authHeaders = @{
    "Authorization" = "Bearer fake_token_1234567890"
}

$createAuthorResponse = Invoke-APIRequest -Method "POST" -Url "$apiUrl/authors" -Body $authorData -Headers $authHeaders
Write-Host "   Status: $($createAuthorResponse.StatusCode)" -ForegroundColor White
if ($createAuthorResponse.Success) {
    Write-Host "   Author ID: $($createAuthorResponse.Content.data.id)" -ForegroundColor Gray
    $authorId = $createAuthorResponse.Content.data.id
} else {
    Write-Host "   Error: $($createAuthorResponse.Content)" -ForegroundColor Red
}
Write-Host ""

# Test 6: Create Category (with auth)
Write-Host "6. Testing Create Category..." -ForegroundColor Cyan
$categoryData = @{
    name = "Fiction"
    description = "Fiction books"
} | ConvertTo-Json

$createCategoryResponse = Invoke-APIRequest -Method "POST" -Url "$apiUrl/categories" -Body $categoryData -Headers $authHeaders
Write-Host "   Status: $($createCategoryResponse.StatusCode)" -ForegroundColor White
if ($createCategoryResponse.Success) {
    Write-Host "   Category ID: $($createCategoryResponse.Content.data.id)" -ForegroundColor Gray
    $categoryId = $createCategoryResponse.Content.data.id
} else {
    Write-Host "   Error: $($createCategoryResponse.Content)" -ForegroundColor Red
}
Write-Host ""

# Test 7: Create Book (with auth)
Write-Host "7. Testing Create Book..." -ForegroundColor Cyan
$bookData = @{
    title = "Test Book"
    isbn = "1234567890123"
    description = "A test book"
    price = 29.99
    stock = 10
    author_id = $authorId
    category_id = $categoryId
} | ConvertTo-Json

$createBookResponse = Invoke-APIRequest -Method "POST" -Url "$apiUrl/books" -Body $bookData -Headers $authHeaders
Write-Host "   Status: $($createBookResponse.StatusCode)" -ForegroundColor White
if ($createBookResponse.Success) {
    Write-Host "   Book ID: $($createBookResponse.Content.data.id)" -ForegroundColor Gray
    $bookId = $createBookResponse.Content.data.id
} else {
    Write-Host "   Error: $($createBookResponse.Content)" -ForegroundColor Red
}
Write-Host ""

# Test 8: Search Books
Write-Host "8. Testing Search Books..." -ForegroundColor Cyan
$searchResponse = Invoke-APIRequest -Method "GET" -Url "$apiUrl/books/search?q=test"
Write-Host "   Status: $($searchResponse.StatusCode)" -ForegroundColor White
Write-Host "   Found: $($searchResponse.Content.data.Count) books" -ForegroundColor Gray
Write-Host ""

# Test 9: Get Books by Author
Write-Host "9. Testing Get Books by Author..." -ForegroundColor Cyan
$booksByAuthorResponse = Invoke-APIRequest -Method "GET" -Url "$apiUrl/books/author/$authorId"
Write-Host "   Status: $($booksByAuthorResponse.StatusCode)" -ForegroundColor White
Write-Host "   Found: $($booksByAuthorResponse.Content.data.Count) books" -ForegroundColor Gray
Write-Host ""

# Test 10: Update Book Stock
Write-Host "10. Testing Update Book Stock..." -ForegroundColor Cyan
$stockData = @{
    stock = 15
} | ConvertTo-Json

$updateStockResponse = Invoke-APIRequest -Method "PUT" -Url "$apiUrl/books/$bookId/stock" -Body $stockData -Headers $authHeaders
Write-Host "   Status: $($updateStockResponse.StatusCode)" -ForegroundColor White
Write-Host "   Message: $($updateStockResponse.Content.message)" -ForegroundColor Gray
Write-Host ""

Write-Host "=== API Testing Complete ===" -ForegroundColor Green
Write-Host "Check the server logs for detailed request information." -ForegroundColor Yellow
