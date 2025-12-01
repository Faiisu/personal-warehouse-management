# Personal Warehouse Management

A full-stack warehouse management system designed to help users track products, categories, stock levels, and warehouse events.

## Features

-   **User Authentication**: Secure registration and login functionality.
-   **Product Management**: Create, read, update, and delete products.
-   **Category Management**: Organize products into categories.
-   **Warehouse Operations**: Manage stock levels and track warehouse movements.
-   **Event Tracking**: Monitor system events and activities.

## Tech Stack

### Frontend
-   **Framework**: React 19
-   **Build Tool**: Vite
-   **Language**: TypeScript
-   **Styling**: CSS (Vanilla)

### Backend
-   **Language**: Go 1.24
-   **Web Framework**: Fiber v2
-   **Database**: MongoDB
-   **Documentation**: Swagger (implied by dependencies)

## Getting Started

### Prerequisites
-   Node.js (v20+ recommended)
-   Go (v1.24+)
-   MongoDB instance

### Installation

#### Backend
1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Run the server:
    ```bash
    go run main.go
    ```
    The server typically runs on port 3000 or 8080 (check `.env` or logs).

#### Frontend
1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the development server:
    ```bash
    npm run dev
    ```

## API Endpoints

### Auth
-   `POST /api/register` - Register a new user
-   `POST /api/login` - Login user

### Products
-   `GET /api/products` - List all products
-   `POST /api/products` - Create a product
-   `PUT /api/products/:productId` - Update a product
-   `DELETE /api/products/:productId` - Delete a product

### Categories
-   `GET /api/categories` - List all categories
-   `POST /api/categories` - Create a category
-   `DELETE /api/categories/:categoryId` - Delete a category

### Warehouse
-   `GET /api/warehouse` - List warehouse stock
-   `POST /api/warehouse` - Add stock
-   `DELETE /api/warehouse/:stockId` - Remove stock

### Health
-   `GET /api/health` - Health check
