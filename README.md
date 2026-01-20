# SiteSentinel

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-gray.svg)](https://opensource.org/licenses/MIT)

SiteSentinel is a high-performance microservice built in Go for monitoring website availability. It tracks HTTP status codes, measures response latency, and dispatches real-time notifications regarding status changes via an asynchronous message broker.

---

## Key Features

* **RESTful API**: Manage monitored endpoints and user-site associations.
* **Automated Health Checks**: Periodic, scheduled probing of registered URLs.
* **Asynchronous Processing**: Utilizes RabbitMQ to decouple monitoring logic from data persistence.
* **Reliable Storage**: PostgreSQL handles site configurations and historical time-series probe data.
* **Clean Architecture**: Deep separation of concerns (Domain, Usecase, Repository, Transport).
* **Structured Logging**: High-performance diagnostics using the native `slog` library.

<p align="center">
<img width="300" alt="arch" src="https://github.com/user-attachments/assets/0c35f88e-e6f2-4d77-8238-a46749d7b340" />
</p>

## Tech Stack

| Category | Technology |
| :--- | :--- |
| **Language** | Go 1.24+ |
| **Web Framework** | [Gin Gonic](https://github.com/gin-gonic/gin) |
| **Message Broker** | [RabbitMQ](https://www.rabbitmq.com/) |
| **Database** | [PostgreSQL](https://www.postgresql.org/) |
| **Dependency Injection** | [Google Wire](https://github.com/google/wire) |
| **Testing** | Testify, Sqlmock, Mockery |
| **Deployment** | Docker & Docker Compose |

---

## API Reference

### Manage Sites

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/addSite` | Register a new URL for monitoring |
| `GET` | `/api/v1/getUserSites/:userID` | Retrieve all sites associated with a specific user |

---

## Database Schema
<p align="center">
    <img width="300" alt="schema" src="https://github.com/user-attachments/assets/60f603fa-b075-4dfa-a76e-dccccfdd7d24" />
</p>

The persistence layer is organized to support high-frequency monitoring data:

* **Users**: Manages user profiles.
* **Sites**: Stores unique URLs and their specific monitoring intervals.
* **User_Sites**: A join table implementing a many-to-many relationship.
* **Check_Results**: A time-series table storing historical probe results (latency, status, timestamp).

---

## Configuration

The application is configured via environment variables. Create a `.env` file in the `deployments` directory.

### Infrastructure Variables

| Variable | Description | Example Value |
| :--- | :--- | :--- |
| **POSTGRES_USER** | Database administrative user | `admin` |
| **POSTGRES_PASSWORD** | Database administrative password | `secure_password` |
| **MONITOR_DB_URL** | DSN for service database access | `postgres://user:pass@db:5432/db` |
| **RABBIT_URL** | AMQP connection string | `amqp://guest:guest@rabbitmq:5672/` |
| **HTTP_PORT** | API listening port | `8080` |

---

## Installation

1.  **Clone the repository**
    ```bash
    git clone [https://github.com/ReilEgor/SiteSentinel.git](https://github.com/ReilEgor/SiteSentinel.git)
    cd SiteSentinel
    ```

2.  **Environment Setup**
    Configure your `.env` variables in the `deployments` folder based on the table above.

3.  **Run with Docker**
    ```bash
    docker-compose up --build
    ```

---

## License

Distributed under the MIT License. See `LICENSE` for more information.

**Developed by [YehorReil](https://github.com/ReilEgor)**
