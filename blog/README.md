# /blog/

This folder contains the Go backend for the `blog-in-golang` project.

## Code Coverage

| File Path                                                 | Function                     | Coverage |
| :-------------------------------------------------------- | :--------------------------- | :------- |
| `github.com/etzelm/blog-in-golang/app.go:22`              | `main`                       | `85.0%`  |
| `github.com/etzelm/blog-in-golang/app.go:55`              | `LoadStaticFileRoutes`       | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:71`              | `LoadServerRoutes`           | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:91`              | `LoadMiddlewares`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:100`             | `staticCacheMiddleware`      | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:117`             | `unauthorizedMiddleware`     | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:133`             | `randRange`                  | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:19` | `AuthPage`                   | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:36` | `SecurePage`                 | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:74` | `AuthResponse`               | `50.0%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:167`| `HashPassword`               | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:172`| `CheckPasswordHash`          | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:22` | `createAWSConfig`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:44` | `createDynamoDBClient`       | `75.0%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:53` | `PostPage`                   | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:73` | `CategoryPage`               | `75.0%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:107`| `ArticlePage`                | `62.5%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:142`| `AboutPage`                  | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:158`| `ContactPage`                | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:179`| `ContactResponse`            | `79.4%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:244`| `renderErrorPage`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:18`| `createS3Client`           | `75.0%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:27`| `ListingsGETAPI`           | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:37`| `ListingGETAPI`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:59`| `ListingPOSTAPI`           | `69.2%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:137`| `UploadImagePOSTAPI`       | `83.3%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:71`     | `createDynamoDBClient`       | `80.0%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:97`     | `GetArticlePanels`           | `87.5%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:172`    | `GetCategoryPageArticlePanels` | `85.4%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:249`    | `GetArticleByID`             | `85.3%`  |
| `github.com/etzelm/blog-in-golang/src/models/realtor.models.go:39`  | `GetRealtorListings`         | `80.8%`  |
| `github.com/etzelm/blog-in-golang/src/models/realtor.models.go:101` | `GetRealtorListing`          | `76.9%`  |
| **Total** |                              | **`82.7%`** |

## Key Components

* **`app.go`**: The main application file that sets up the Gin web server, defines routes, and loads middleware. It handles serving static files, blog post pages, category pages, individual articles, and the contact page. It also includes logic for CertMagic for automatic HTTPS in production environments.
* **`docker-compose.yml`**: Defines the Docker service for local development, specifying the image, container name, ports, and environment variables (like AWS credentials and the DynamoDB table name for articles).
* **`Dockerfile`**: A multi-stage Dockerfile that:
    1. Builds the Go application.
    2. Builds the React frontend application (from the `/realtor` folder).
    3. Creates a final minimal runtime image containing the Go binary and the built React assets. It also sets up a non-root user for security.
* **`go.mod`**: Defines the Go module and its dependencies, including Gin, AWS SDK, CertMagic, and Logrus.
* **`/src/`**:
  * **`/handlers/`**: Contains the Go HTTP handlers for different routes:
    * `blog.handlers.go`: Handles requests related to blog posts, categories, individual articles, about page, and contact form submissions.
    * `auth.handlers.go`: Manages user authentication, including displaying an auth page and handling login/secure page access (with bcrypt for password hashing).
    * `realtor.handlers.go`: Provides API endpoints for the realtor frontend, including fetching all listings, a specific listing, adding/updating listings in DynamoDB, and uploading images to S3.
  * **`/models/`**: Defines the data structures (structs) used in the application:
    * `blog.models.go`: Defines `ContactForm`, `Item` (raw DynamoDB article structure), `Article` (processed article structure with `template.HTML`), and `Category`. Includes functions to fetch article panels and individual articles from DynamoDB.
    * `realtor.models.go`: Defines the `Listing` struct for real estate properties and includes functions to get all listings or a specific listing from DynamoDB.
    * `auth.models.go`: Defines the `AuthForm` struct for authentication.
* **`/templates/`**: (Assumed based on `httpServer.LoadHTMLGlob("templates/*")` in `app.go`) Contains Go HTML templates used for rendering the blog's frontend (e.g., `index.html`, `article.html`, `contact.html`, `about.html`, `error.html`, `auth.html`, `secure.html`).
* **`/public/`**: (Assumed based on `LoadStaticFileRoutes` in `app.go`) Contains static assets like `robots.txt`, `sitemap.xml`, images (`favicon.ico`), and potentially CSS/JS for the blog's non-React parts.
* **`app_test.go`**: Contains unit tests for the `app.go` functionalities, including testing random number generation, middleware (static cache, unauthorized access), route loading, and main execution paths (with and without CertMagic). It utilizes standard Go testing, `httptest` for HTTP requests, and mocks/stubs where necessary.

## Functionality

* Serves a blog with articles stored in AWS DynamoDB.
* Provides API endpoints for a React-based realtor listing application.
* Supports user authentication via a custom Go implementation.
* Handles image uploads to AWS S3.
* Uses Gin middleware for gzip compression and caching.
* Includes security middleware to block common malicious request paths.
* Designed for deployment via Docker.
