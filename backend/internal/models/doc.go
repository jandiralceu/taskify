/*
Package models defines the core domain entities used for business logic and database persistence.

Core entities include:
  - User: The main user profile (including credentials, activity, and avatar metrics).
  - Task: The primary piece of work, with priority, status, and assignment logic.
  - TaskNote & TaskAttachment: Child entities for extended task collaboration.
  - RoleModel & Permission: RBAC-related data representations for roles-to-permissions mappings.

These structs are mapped to PostgreSQL tables via GORM's database tags.
*/
package models
