<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->




<!-- PROJECT LOGO -->
<br />



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#docs">Documentation</a></li>
    <li><a href="#db">DB Diagram</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project


the backend of CollabTED

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Next][Golang]][Golang-url]
* [![React][Echo]][Echo-url]
* [![Vue][Prisma]][Prisma-url]
* [![Angular][GraphQl]][GraphQl-url]
* [![Svelte][Redis]][Redis-url]
* [![Laravel][Docker]][Docker-url]
* [![Bootstrap][Nginx]][Nginx-url]
* [![JQuery][MongoDB]][MongoDB-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites
you need to have go installed

# Installation
## Normal Setup
1. Clone the repo 

    ```sh
   git clone https://github.com/CollabTED/CollabTed-Backend
   ```
2. Generate the db

    ```sh
   go run github.com/steebchen/prisma-client-go generate --schema ./prisma
   go run github.com/steebchen/prisma-client-go db push --schema ./prisma
   ```
3. Install the required packages

   ```sh
   go mod tidy
   ```
5. Edit your .env
6. Runt the project
   ```sh
   go run cmd/server/main.go
   ```
## Using Docker
```sh
docker-compose build app
docker-compose build web
docker-compose up web -d
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Documentation
Check the docs at /swagger/index.html


<!-- USAGE EXAMPLES -->

# DB Diagram



<!-- ROADMAP -->
# Roadmap

- [x] Jwt Auth
- [x] Google And Facebook Oath
- [x] Email Verification
- [x] Profile Managmenent
- [x] Notifications with SSE

<p align="right">(<a href="#readme-top">back to top</a>)</p>













<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[product-screenshot]: https://i.imgur.com/Bd51sZB.png
[Golang]:https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
[Echo]: https://img.shields.io/badge/echo-35495E?style=for-the-badge&logo=gin&logoColor=white
[Echo-url]:https://echo.labstack.com
[Prisma]: https://img.shields.io/badge/Prisma-3982CE?style=for-the-badge&logo=Prisma&logoColor=white
[Prisma-url]: https://prisma.io
[GraphQl]: https://img.shields.io/badge/-GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white
[GraphQl-url]: https://graphql.org/
[Redis]: https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white
[Redis-url]: https://redis.io/
[Docker]: https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://docker.com
[Nginx]: https://img.shields.io/badge/nginx-%23009639.svg?style=for-the-badge&logo=nginx&logoColor=white
[Nginx-url]: https://nginx.org/ 
[MongoDB]: https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white
[MongoDB-url]: https://mongodb.com/

