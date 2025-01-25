# Los Complejos Backend

Welcome to the **Los Complejos Backend** project! This repository is a Go-based backend application designed for managing users (Complejos) and events efficiently. The application utilizes MongoDB as the database and implements JWT-based authentication to ensure security.

---

## 📋 Features

- **User Management**: Create, update, and manage user profiles with optional fitness-related attributes (weight, height, bench, squat, deadlift).
- **Event Management**: Admin-only event creation, user subscription, and unsubscription functionality.
- **JWT Authentication**: Secure access to endpoints using JSON Web Tokens.
- **Role-Based Access Control**: Differentiate between `admin` and `user` roles for controlled access to features.
- **MongoDB Integration**: High-performance database operations with MongoDB.

---

## ⚙️ Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/intxz/los-complejos-backend.git
   cd los-complejos-backend
   ```

2. **Set Up Environment Variables**:
   Create a `.env` file in the root directory and add your MongoDB URI and JWT secret key:
   ```plaintext
   MONGO_URI=mongodb://localhost:27017
   JWT_SECRET=your_secret_key
   ```

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the Application**:
   ```bash
   go run main.go
   ```
   The server will start on [http://localhost:8080](http://localhost:8080).

---

## 🛠️ API Endpoints

### **Authentication**
JWT-based authentication using the `Authorization` header.

### **User (Complejo) Management**

| Method | Endpoint          | Description                       |
|--------|-------------------|-----------------------------------|
| POST   | `/complejo`       | Create a new user (Complejo).     |
| GET    | `/complejo`       | Retrieve all users.               |
| GET    | `/complejo/:id`   | Retrieve a specific user by ID.   |
| PUT    | `/complejo/admin` | Update any user (Admin only).     |
| PUT    | `/complejo/user`  | Update self (User role only).     |

### **Event Management**

| Method | Endpoint                    | Description                          |
|--------|-----------------------------|--------------------------------------|
| POST   | `/event`                    | Create a new event (Admin only).     |
| GET    | `/event`                    | Retrieve all events.                 |
| GET    | `/event/:id`                | Retrieve a specific event by ID.     |
| PUT    | `/event/:id/subscribe`      | Subscribe to an event.               |
| PUT    | `/event/:id/unsubscribe`    | Unsubscribe from an event.           |

---

## 📂 Project Structure

```
los-complejos-backend/
│
├── database/          # MongoDB connection and utilities
├── handlers/          # API endpoint handlers
├── middleware/        # Authentication and authorization middleware
├── models/            # Data models for users (Complejo) and events
├── utils/             # Utility functions (e.g., JWT, IMC calculation)
├── .env               # Environment variables (not tracked by Git)
├── go.mod             # Go module dependencies
├── main.go            # Entry point of the application
└── README.md          # Project documentation
```

---

## ✨ Key Highlights

- **IMC Classification**: Calculate and classify users into fun categories like "NPC" and "Burger King Slayer" based on their fitness metrics.
- **Admin-Only Features**: Event creation and unrestricted user updates are limited to admins.
- **Subscription System**: Users can subscribe or unsubscribe from events, with proper conflict handling.

---

## 🚀 Future Improvements

- Add pagination and filtering for user and event queries.
- Enhance security by implementing rate-limiting.
- Integrate advanced error handling and logging.
- Expand test coverage with unit and integration tests.

---

## 🤝 Contributing

Contributions are welcome! Feel free to fork this repository, submit pull requests, or suggest features via issues.

---

## 👨‍💻 Developed by intxz (Ángel Redondo Gamero)
