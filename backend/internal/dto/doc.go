/*
Package dto defines the Data Transfer Objects used to decouple internal models from the public API schema.

By using distinct structs for requests (binding) and responses (JSON tags), the application can:
  - Version the API independently of the database.
  - Enforce strict validation rules (e.g., email format, range checks).
  - Selectively hide sensitive fields (e.g., password hashes) without using internal model tags.
  - Standardize generic structures like PaginatedResponse[T].
*/
package dto
