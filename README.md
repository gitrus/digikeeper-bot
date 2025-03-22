# digikeeper-bot
Digikeeper bot is based on golang library `telego`.
It was created as pet-project to demonstrate session and stetfull-protocol based approaches in telegram bot development.
To shocase the differences, there are two types of interactions:
- `/command-s` (similar to HTTP-handlers)
- keyboars actions (FSM + Session related)

Current features that you can easly copy-paste to your bot:
- `/cancel` command to erase current state in session
- `/help` command is autogenerated
- `UserStateMiddleware` -- middleware to fetch session by userID (with inmemory mock implementation)

## Docs
To deep dive into WHY/HOW this bot is developed, please refer to:
`./docs/` -- architecture, design, etc.
`README.md` -- overall description

This documentation is central part of the project, not just an addendum. Docs is essential for understanding the project's structure and goal.

## Run
`Dockerfile` -- builded and located as package on github

`./deployments/docker-compose.yml` -- to run it

OR:
run `go run cmd/bot` in root directory

## Contribution
Open a PR or Issue, I will be glad to see it.

[*] Common conduct rules are applied (github-repo related contribution rules).
[*] Apache 2.0 license is for safe usage and sharing

Coding standarts:
- `go tool` is used
