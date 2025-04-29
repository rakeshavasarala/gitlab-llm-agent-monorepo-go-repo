build:
	cd agents/approval-agent && go build -o ../../bin/approval-agent
	cd agents/security-labeler && go build -o ../../bin/security-labeler
	cd agents/reasoning-labeler && go build -o ../../bin/reasoning-labeler

run-approval:
	cd agents/approval-agent && go run main.go

run-security:
	cd agents/security-labeler && go run main.go

run-reasoning:
	cd agents/reasoning-labeler && go run main.go
