/*
Package service handles the application's central business logic and task orchestration.

Services are responsible for:
  - Performing multi-step business transactions.
  - Hashing passwords via Argon2 before persistence.
  - Managing file storage lifecycle for avatars and attachments.
  - Coordinating between multiple repositories (e.g., retrieving User and Roles).
  - Enforcing logic invariants before delegating to the persistence layer.
*/
package service
