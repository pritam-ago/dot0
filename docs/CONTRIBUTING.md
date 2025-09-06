# Contributing to CloudStore

Thank you for your interest in contributing to CloudStore! This document provides guidelines and instructions for contributing to the project.

## 🤝 Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## 🚀 Getting Started

1. **Fork the Repository**

   - Click the Fork button in the top right of the repository
   - Clone your fork locally

2. **Set Up Development Environment**

   ```bash
   # Clone your fork
   git clone https://github.com/your-username/cloudstore.git
   cd cloudstore

   # Set up development environment
   ./scripts/setup.sh
   ```

## 💻 Development Workflow

1. **Create a Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**

   - Write your code
   - Follow the coding standards
   - Add tests for new features

3. **Test Your Changes**

   ```bash
   ./scripts/test.sh
   ```

4. **Lint Your Code**

   ```bash
   ./scripts/lint.sh
   ```

5. **Commit Your Changes**

   ```bash
   # Follow conventional commits
   git commit -m "feat: add new feature"
   ```

6. **Push to Your Fork**

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**
   - Go to the repository on GitHub
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template

## 📝 Coding Standards

### Go

- Follow the official Go style guide
- Use `gofmt` for formatting
- Add comments for exported functions
- Write unit tests

### TypeScript/React

- Use TypeScript for type safety
- Follow ESLint configuration
- Use functional components
- Write unit tests with Jest

### Git Commits

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation
- `style:` Formatting
- `refactor:` Code restructuring
- `test:` Tests
- `chore:` Maintenance

## 🧪 Testing

1. **Unit Tests**

   - Write tests for new features
   - Maintain existing tests
   - Aim for good coverage

2. **Integration Tests**

   - Test component interactions
   - Test API endpoints
   - Test database operations

3. **End-to-End Tests**
   - Test complete workflows
   - Test UI interactions

## 📚 Documentation

1. **Code Documentation**

   - Document public APIs
   - Add inline comments
   - Update README files

2. **Architecture Documentation**
   - Update diagrams
   - Document new features
   - Explain design decisions

## 🐛 Bug Reports

When reporting bugs:

1. Check existing issues
2. Use the bug report template
3. Include:
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment details

## 💡 Feature Requests

When suggesting features:

1. Check existing issues
2. Use the feature request template
3. Explain the use case
4. Describe expected behavior

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.
