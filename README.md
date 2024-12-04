# Go-Lox Interpreter 🚀

An implementation of the Lox programming language Interpreter written in Go. This interpreter is based on the book ["Crafting Interpreters"](https://craftinginterpreters.com/) by Robert Nystrom and tested via [Codecrafters](https://app.codecrafters.io/courses/interpreter).

## About Lox 📖

Lox is a dynamically-typed scripting language that supports object-oriented programming with classes and inheritance. This implementation follows the tree-walk interpreter pattern.

## Features ✨

Here's what's currently implemented in this interpreter:

| Feature | Status |
|---------|---------|
| Basic Arithmetic | ✅ |
| Variables | ✅ |
| Control Flow (if/else) | ✅ |
| Loops (while, for) | ✅ |
| Functions | ✅ |
| Closures | ✅ |
| Classes | 🚧 |
| Inheritance | 🚧 |
| Standard Library | 🚧 |
| Error Handling | 🚧 |

## Getting Started 🛠️

### Prerequisites

- Go 1.22 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-lox.git

# Navigate to project directory
cd go-lox
```

### Running the Interpreter

There are four ways to use the interpreter:

1. Print tokens of a Lox script:
```bash
./go-lox tokenize <filename>.lox
```

2. Print AST of a Lox script:
```bash
./go-lox parse <filename>.lox
```

3. Evaluate basic expression:
```bash
./go-lox eval <filename>.lox
```

4. Run a Lox script:
```bash
./go-lox run <filename>.lox
```

## Example Lox Program 📝

```lox
// This is a simple Lox program
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

print fib(10); // Outputs: 55
```

## Project Structure 🏗️

```
├── parser/       # Abstract Syntax Tree parser
├── lexer/        # Lexical analysis
├── interpreter/  # Interpreter implementation
├── expression/   # Expression definitions e.g., <, ==, +, >, -
├── stmt/         # Statement definitions e.g., var, fun, for, while
├── token/        # Token definition
└── main.go       # Entry point
```

## Contributing 🤝

I will be glad to receive any of your questions/suggestions/contributions to this project! Feel free to open a PR or contact me via:

[Twitter](https://x.com/4c656f)

[Email](mailto:tarabrinleonid@gmail.com)

[Telegram](https://t.me/c656f)

---

For more information about the Lox language, visit [Crafting Interpreters](https://craftinginterpreters.com/).

*Note: This implementation is for educational purposes and may not be suitable for production use.*