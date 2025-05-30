# /blog/

This folder contains the Go backend for the `blog-in-golang` project.

## Code Coverage

| File Path                                                 | Function                     | Coverage |
| :-------------------------------------------------------- | :--------------------------- | :------- |
| `github.com/etzelm/blog-in-golang/app.go:22`              | `main`                       | `85.0%`  |
| `github.com/etzelm/blog-in-golang/app.go:56`              | `LoadStaticFileRoutes`       | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:72`              | `LoadServerRoutes`           | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:92`              | `LoadMiddlewares`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:101`             | `staticCacheMiddleware`      | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:121`             | `unauthorizedMiddleware`     | `100.0%` |
| `github.com/etzelm/blog-in-golang/app.go:140`             | `randRange`                  | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:20` | `AuthPage`                   | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:37` | `SecurePage`                 | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:75` | `AuthResponse`               | `81.8%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:143`| `HashPassword`               | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/auth.handlers.go:148`| `CheckPasswordHash`          | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:21` | `PostPage`                   | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:41` | `CategoryPage`               | `87.5%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:75` | `ArticlePage`                | `87.5%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:110`| `AboutPage`                  | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:126`| `ContactPage`                | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:147`| `ContactResponse`            | `96.6%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/blog.handlers.go:211`| `renderErrorPage`            | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:21`| `ListingsGETAPI`           | `100.0%` |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:31`| `ListingGETAPI`            | `66.7%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:53`| `ListingPOSTAPI`           | `90.5%`  |
| `github.com/etzelm/blog-in-golang/src/handlers/realtor.handlers.go:116`| `UploadImagePOSTAPI`       | `88.9%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:69`     | `GetArticlePanels`           | `91.2%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:138`    | `GetCategoryPageArticlePanels` | `91.4%`  |
| `github.com/etzelm/blog-in-golang/src/models/blog.models.go:209`    | `GetArticleByID`             | `96.6%`  |
| `github.com/etzelm/blog-in-golang/src/models/realtor.models.go:41`  | `GetRealtorListings`         | `86.4%`  |
| `github.com/etzelm/blog-in-golang/src/models/realtor.models.go:100` | `GetRealtorListing`          | `81.8%`  |
| **Total** |                              | **`91.2%`** |

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
