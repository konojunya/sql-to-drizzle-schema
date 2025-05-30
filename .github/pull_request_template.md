# Pull Request

## Description
<!-- Provide a brief description of your changes -->

## Type of Change
<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📝 Documentation update
- [ ] 🧪 Test improvements
- [ ] 🔧 Code refactoring
- [ ] ⚡ Performance improvement
- [ ] 🔒 Security improvement

## Changes Made
<!-- List the specific changes made in this PR -->

- 
- 
- 

## Testing
<!-- Describe the tests you ran to verify your changes -->

- [ ] I have run `make test` and all tests pass
- [ ] I have run `make check` and all linting/formatting checks pass
- [ ] I have added/updated tests for my changes
- [ ] I have tested the changes manually

## SQL Dialect Support
<!-- If applicable, mark which SQL dialects were tested -->

- [ ] PostgreSQL
- [ ] MySQL (if applicable)
- [ ] Spanner (if applicable)

## Example Usage
<!-- If applicable, provide an example of how to use the new feature -->

```bash
# Example command
./sql-to-drizzle-schema example.sql -o output.ts
```

```sql
-- Example SQL input
CREATE TABLE example (...);
```

```typescript
// Expected TypeScript output
export const example = pgTable('example', {...});
```

## Documentation
<!-- Check if documentation updates are needed -->

- [ ] I have updated the README.md if needed
- [ ] I have updated CLAUDE.md if needed
- [ ] I have added/updated code comments
- [ ] No documentation changes required

## Checklist
<!-- Ensure your PR meets these requirements -->

- [ ] My code follows the project's coding style
- [ ] I have performed a self-review of my own code
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Related Issues
<!-- Link any related issues -->

Closes #<!-- issue number -->
Relates to #<!-- issue number -->

## Screenshots (if applicable)
<!-- Add screenshots to help explain your changes -->

## Additional Notes
<!-- Add any additional notes or context about the PR -->

---

**Reviewer Notes:**
<!-- For reviewers: Any specific areas to focus on during review -->
