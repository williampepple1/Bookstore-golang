# Database Migrations

This directory contains SQL migration files for the bookstore database. The migration system supports both manual SQL migrations and GORM auto-migrations.

## Migration Files

Migration files should be named with a sequential number prefix followed by a descriptive name:
- `000_init.sql` - Database initialization
- `001_create_authors_table.sql` - Create authors table
- `002_create_categories_table.sql` - Create categories table
- `003_create_books_table.sql` - Create books table
- `004_add_book_ratings_table.sql` - Add book ratings table

## Running Migrations

### Using Make commands:
```bash
# Run all pending migrations
make migrate

# Check migration status
make migrate-status

# Rollback last migration
make migrate-rollback

# Validate migration files
make migrate-validate
```

### Using the migrate CLI tool directly:
```bash
# Run migrations
go run cmd/migrate/main.go -action=migrate

# Check status
go run cmd/migrate/main.go -action=status

# Rollback
go run cmd/migrate/main.go -action=rollback

# Validate
go run cmd/migrate/main.go -action=validate
```

## Migration System Features

1. **Automatic Migration Tracking**: The system automatically tracks which migrations have been applied using a `migrations` table.

2. **Transaction Safety**: Each migration runs in a transaction, so if it fails, the database is rolled back to its previous state.

3. **Idempotent**: Migrations can be run multiple times safely - already applied migrations are skipped.

4. **Validation**: Migration files are validated before execution to ensure they're properly formatted.

5. **Status Reporting**: You can check which migrations have been applied and when.

6. **Rollback Support**: Basic rollback functionality is available (removes migration records).

## Best Practices

1. **Always use IF NOT EXISTS**: Use `CREATE TABLE IF NOT EXISTS` and `CREATE INDEX IF NOT EXISTS` to make migrations idempotent.

2. **Use transactions**: Wrap related operations in transactions when possible.

3. **Test migrations**: Always test migrations on a copy of your production data before applying to production.

4. **Backup before major changes**: Always backup your database before running migrations that modify existing data.

5. **Use descriptive names**: Migration file names should clearly describe what they do.

6. **Version control**: Always commit migration files to version control.

## Migration File Structure

Each migration file should:
- Be a valid SQL file
- Use PostgreSQL syntax
- Include proper error handling with `IF NOT EXISTS`
- Create necessary indexes for performance
- Include foreign key constraints where appropriate
- Use the `update_updated_at_column()` function for timestamp updates

## Troubleshooting

If a migration fails:
1. Check the error message in the logs
2. Fix the SQL in the migration file
3. The migration will be retried on the next run
4. If needed, you can rollback using `make migrate-rollback`

For more complex rollbacks, you may need to manually modify the database or create a new migration to fix issues.
