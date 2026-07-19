# 🌍 Smart Waste Management API & Routing Engine

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-20232A?style=for-the-badge&logo=go&logoColor=white)
![Leaflet](https://img.shields.io/badge/Leaflet-199900?style=for-the-badge&logo=leaflet&logoColor=white)

## About the Project

Smart Waste Management is a high-performance backend service that transforms urban waste management from static routes into a dynamic, data-driven, and emergency-alert system. 

Using telemetry data (fill level, temperature, methane gas level) received from IoT sensors, the algorithm calculates the most optimal and safe collection route for garbage trucks in real-time.

## ⚙️ Core Features and Architectural Decisions

*   **Smart Route Optimization (SQL-Driven):** Route calculations are not done on the application layer, but directly on PostgreSQL using weighted scoring and conditional logic for high performance.
*   **Emergency Circuit Breaker:** Even if the fill rate is only 10%, if the temperature inside the container exceeds 65°C or the Methane Gas (CH4) level exceeds 20 ppm, the system overrides the physical distance sorting and elevates that container to the 1st priority (Fire and Explosion prevention).
*   **ACID Compliant Data Consistency:** The processes of logging waste and updating the container's fill level are protected by **SQL Transactions** in Go. If either process fails, a full Rollback is triggered to ensure data integrity.
*   **Carbon Footprint Calculation:** On the domain layer, instant carbon footprint logging is calculated using emission multipliers based on the waste type and weight.
*   **Live Navigation and Simulation:** Using OSRM (Open Source Routing Machine) and Leaflet.js integration, the mathematically calculated sequence is mapped onto the real-world street network and visualized with directional arrows.

## 🚀 Quick Start

### Prerequisites
*   Go (1.20+)
*   PostgreSQL
*   Git

### Installation

1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd smartwaste

### Set up environment variables:
Create a .env file in the root directory and enter your database credentials:

DB_URL=postgres://username:password@localhost:5432/smart_waste
PORT=3000
Download dependencies and start the server:

###
go mod tidy
go run cmd/api/main.go

### Running the simulator 
go run cmd/simulator/main.go