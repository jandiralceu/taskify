/*
Package repository manages the persistence and retrieval of data from the PostgreSQL database.

Repositories are abstraction layers over the GORM-DB connection, focusing on:
  - CRUD operations (Create, Read, Update, Delete) on core entities.
  - SQL query construction with parameterization (SQL injection protection).
  - Advanced filtering, searching, and pagination logic.
  - Preloading and joining associations (e.g., loading Roles for Users).
*/
package repository
