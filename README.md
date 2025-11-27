# 🇦🇺 Australia Postcode Scraper API (Go/Docker)

This project provides a simple, high-performance RESTful API service
written in **Go (Golang)** to scrape and return Australian postcode and
suburb data based on a keyword search from the **Australia Post
website**.

## ✨ Author
-   Name: Ryan Zeng
-   Email: zengruijiang@gmail.com

## ✨ Features

-   **RESTful Endpoint** -- `/search` endpoint to query postcode data
-   **Data Extraction** -- Scrapes postcode, suburb name, and state
    (e.g., `NSW`)
-   **JSON Output** -- Clean, structured JSON responses
-   **Dockerized** -- Lightweight and ready for container deployment

## ⚙️ Requirements

To run this application, you need:

-   **Go 1.21+** (for local development)
-   **Docker** (for containerized deployment)

## 🏗️ Project Structure

The core files required are:

-   `postcode_scraper.go` --- main Go application with HTTP server +
    scraping logic\
-   `go.mod` / `go.sum` --- Go module and dependency files\
-   `Dockerfile` --- instructions for container build

## 🚀 Running the Application

### 1. Local Development (Without Docker)

#### Initialize Go Module & Install Dependencies

``` bash
go mod init postcode_scraper
go get github.com/PuerkitoBio/goquery
go mod tidy
```

#### Run the Server

``` bash
go run postcode_scraper.go
```

#### Test the Endpoint

``` bash
curl http://localhost:8080/search?keyword=sydney
```

### 2. Dockerized Setup (Recommended)

#### Ensure Dependencies Are Ready

``` bash
go mod init postcode_scraper
go mod tidy
```

#### Build Image

``` bash
docker build -t postcode-api .
```

#### Run Container

``` bash
docker run -d -p 8080:8080 --name postcode-scraper-app postcode-api
```

#### Test Endpoint

``` bash
curl http://localhost:8080/search?keyword=melbourne
```

## 📝 API Usage

### Endpoint

    GET /search

### Query Parameters

  -----------------------------------------------------------------------
  Parameter        Required         Description           Example
  ---------------- ---------------- --------------------- ---------------
  `keyword`        Yes              The suburb or town    `sydney`,
                                    name to search        `brisbane`

  -----------------------------------------------------------------------

### Success Response Example

``` json
[
    {
        "postcode": "2055",
        "suburb": "NORTH SYDNEY",
        "state": "NSW",
        "category": "Delivery Area"
    }
]
```

### Error Responses

``` json
{
    "error": "Missing 'keyword' parameter"
}
```
