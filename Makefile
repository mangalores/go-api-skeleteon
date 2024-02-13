

.PHONY: swag
swag:
	swag init

.PHONY: test
test:
	go clean -testcache
	go vet ./...
	go test ./... -race

.PHONY: cover
cover:
	go test -coverprofile=./docs/coverreport/coverage.out ./...

.PHONY: cover-report
cover-report: cover
	go tool cover -html ./docs/coverreport/coverage.out -o ./docs/coverreport/coverage-report.html

.PHONY: full-report
full-report:
	go clean -testcache
	goreporter -p ./ -r ./docs/staticreport -f html -e vendor -stderrthreshold 1


.PHONY: mocks
mocks:
	mockgen -source ./internal/svc/types.go -destination ./internal/mocks/svc/types.go
	mockgen -source ./internal/svc/daily.go -destination ./internal/mocks/svc/daily.go
	mockgen -source ./internal/svc/derive.go -destination ./internal/mocks/svc/derive.go
	mockgen -source ./internal/svc/historic.go -destination ./internal/mocks/svc/historic.go


