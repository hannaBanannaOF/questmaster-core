# Tasks

## 1. Counts

- [ ] 1.1 Add a repository query that groups the requester's campaigns by status with the `role` filter of the campaign list; verify with repository tests (CI Postgres) for each role, both roles and another user's campaigns
- [ ] 1.2 Add the use case that fills every status with 0 before applying the query result; verify with a unit test for a user with no campaigns

## 2. Transport

- [ ] 2.1 Add `GET /core/api/v1/campaign/counts` with its handler, response and Swagger annotation, registered before `/campaign/{campaignID}`; verify the route with a handler test, including 400 for an unknown `role`

## Workflow follow-up

- Implement after `paginate-list-endpoints`, which adds the `role` filter.
- Archive the change after CI passes.
