# Youtube Explorer
## Project Goal

To make an API to fetch the latest videos sorted in reverse chronological order of their publishing date-time from YouTube for a given tag/search query in a paginated response.

## Features

- YouTube API continuously in the background for fetching the latest videos for a predefined search query and should store the data of videos in a database (here, MongoDB)
- A GET API that returns the stored video data in a paginated response sorted in descending order of published datetime.
- A basic search API to search the stored videos using their title and description.

## Table of Contents
* [Getting Started](##-Getting-Started)
* [Running the Project](##-Running-the-Project)
  - [Run manually](###-Run-Manually)
  - [Run with Docker](###-Run-with-Docker)
* [API Reference](##-API-Reference)

## Getting Started

To get started with this project, follow the steps below:

### Clone the Repository
* Clone the Git repository to your local machine:
```bash
  git clone https://github.com/megsingh/youtube-explorer.git
```
* Navigate to the project directory:
```bash
  cd youtube-explorer
```
### Set Up Environment Variables
* Copy the provided [.env.sample](https://github.com/megsingh/youtube-explorer/blob/main/.env.sample) file to a new file named .env:
```bash
  cp .env.sample .env
```
Edit the .env file and fill in the required values. You may refer to [.env.default](https://github.com/megsingh/youtube-explorer/blob/main/.env.default) for default values.

### Get YouTube Data API Keys
* Visit the [YouTube Data API](https://developers.google.com/youtube/v3/getting-started) page to obtain API keys. Add the multiple API keys to the .env file in the format:

```bash
API_KEY= 'YOUR_API_KEY_1 | YOUR_API_KEY_2 | YOUR_API_KEY_3'
```
## Running the Project

### Run Manually
* Install the project dependencies using your preferred package manager. For example, with Go modules:
```bash
go mod download
```

* Run the project:
```bash
go run main.go
```

### Run with Docker
Make appropriate changes to the environment variables in [compose.yaml](https://github.com/megsingh/youtube-explorer/blob/main/compose.yaml) file and use the following command to build and start the app:

```bash
docker-compose up --build
```

## API Reference

#### Get Latest Videos

```http
  GET /v1/api/videos
```
Retrieve the latest videos sorted in reverse chronological order of their publishing date-time 

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `next_token` | `string` | **Optional**. Pagination Token received from the result of this API call to fetch the next page. If nothing is passed then it returns the first page of results |

#### Get item

```http
  GET /v1/api/search
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `next_token`      | `string` | **Optional**. Pagination Token received from the result of this API call to fetch the next page. If nothing is passed then it returns the first page of results|
| `search_query`      | `string` | **Optional**. tag/search query |

