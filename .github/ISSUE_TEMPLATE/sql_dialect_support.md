---
name: SQL Dialect Support
about: Request support for a specific SQL dialect or feature
title: '[DIALECT] '
labels: ['enhancement', 'sql-dialect', 'needs-triage']
assignees: ['konojunya']
---

# SQL Dialect Support Request

## SQL Dialect
<!-- Which SQL dialect needs support? -->

- [ ] MySQL
- [ ] Spanner  
- [ ] SQLite
- [ ] MariaDB
- [ ] Oracle
- [ ] SQL Server
- [ ] Other: 

## Specific Feature/Syntax
<!-- What specific SQL feature or syntax needs support? -->

### SQL Example
<!-- Provide an example of the SQL syntax that should be supported -->

```sql
-- Example SQL syntax that needs support
CREATE TABLE example (
    -- Your SQL example here
);
```

### Expected Drizzle Output
<!-- What should the generated Drizzle ORM code look like? -->

```typescript
// Expected TypeScript/Drizzle output
export const example = pgTable('example', {
    // Expected schema definition
});
```

## Current Behavior
<!-- What happens when you try to parse this SQL currently? -->

- [ ] Parsing fails with error
- [ ] Parsing succeeds but output is incorrect
- [ ] Feature is completely unsupported
- [ ] Other: 

## Error Messages (if any)
<!-- Include any error messages you received -->

```
Error output here
```

## Documentation References
<!-- Links to official documentation for this SQL feature -->

- [ ] Official SQL dialect documentation: 
- [ ] Drizzle ORM documentation: 
- [ ] Other references: 

## Use Case Priority
<!-- How common is this SQL feature? -->

- [ ] Very common - used in most schemas
- [ ] Common - used frequently  
- [ ] Uncommon - used occasionally
- [ ] Rare - used in specific cases

## Compatibility Considerations
<!-- Are there any compatibility issues to consider? -->

## Implementation Complexity
<!-- Any thoughts on implementation complexity? -->

- [ ] Simple - straightforward mapping
- [ ] Medium - requires some logic
- [ ] Complex - significant parser changes needed
- [ ] Very Complex - major architectural changes

## Additional Context
<!-- Add any other context about this SQL dialect support request -->

## Related Issues
<!-- Link any related issues -->

## Sample Schema
<!-- If you have a complete sample schema using this feature, please include it -->