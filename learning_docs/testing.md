### Table driven testing
Table-driven testing in Go is a testing approach where test cases are defined using a `slice` of *anonymous* structs, allowing multiple test scenarios to be defined concisely. Anonymous structs are particularly useful in this approach as they can quickly define test inputs and expected outputs without creating separate named types. This is analogous to `pytest parametrize` in python.

Meaningful testing focuses on gaining actual value from tests, rather than simply writing tests to demonstrate test coverage. The goal is to validate the functionality and integrity of code, ensuring that tests provide genuine insights into the system's behavior.

Test-Driven Development (TDD) is an approach where tests are written before the actual code implementation. Potential drawbacks include the repetitive nature of running tests multiple times and the potential overhead of constantly switching between writing tests and implementation code.


### Testing library
A popular testing library in Go is [testify](https://github.com/stretchr/testify). Testify is an industry-standard Go testing package that provides utilities for asserting values, checking for errors, and creating mocks in test cases.

Use a separate test database that is wiped clean between test suites to prevent contamination between different test cases, ensuring each test starts with a blank slate.


### Running tests
All test files should end with test.go e.g. `<filename>_test.go`. Within the test file, all functions that need to be tested should start with a test prefix e.g.
```go
func TestCreateWorkout(t *testing.T) {
    ...
}
```

Navigate to the directory with the test files and run `go test . -v`. The verbose flag `-v` is optional to see the steps. By default, you need to be in the directory of the test file to use `go test`. However, you can use a wildcard `./...` to run tests across all packages `go test ./... -v`.


In Go, every folder is a package. Go treats packages as independent units of compilation.
- **Isolation**: If you are working on the workout package, you usually only want to run tests for that package. Running 500 tests across 50 other modules every time you save a file would be slow.
- **Caching**: Go caches test results per package. If you run go test ./... and you only changed code in the workout folder, Go will use the cached results for the other 49 packages, making the total run time nearly instant.
- **Compilation**: go test actually compiles a separate "test binary" for every single package it finds.


When running tests across your whole module, these flags are your best friends:
| Flag | Purpose | 
| :--- | :------ | 
| `-v`  | **Verbose**: Shows the name of every test being run. | 
| `-count=1` | **Force Run**: Bypasses the cache (useful if you're debugging DB issues). | 
| `-parallel 4` | Run tests for different packages at the same time. |
| `-failfast` | Stop the whole suite as soon as a single test fails. | 