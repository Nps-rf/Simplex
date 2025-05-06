# Вклад в проект Simplex

Спасибо, что хотите внести вклад в развитие Simplex!

## Как внести изменения

1. **Форкните репозиторий** и создайте отдельную ветку для вашей задачи:
   ```bash
   git checkout -b feature/ваша-фича
   ```
2. **Пишите чистый, читаемый код**. Следуйте стилю Go и best practices.
3. **Покрывайте изменения тестами**. Покрытие должно быть не ниже 80% для новых модулей.
4. **Запускайте линтер и тесты**:
   ```bash
   golangci-lint run
   go test ./...
   ```
5. **Оформляйте коммиты осмысленно** (на русском или английском):
   - feat: новая функциональность
   - fix: исправление бага
   - docs: документация
   - test: тесты
   - refactor: рефакторинг
6. **Оформляйте Pull Request**:
   - Кратко опишите изменения
   - Укажите связанные issue (если есть)
   - Приложите примеры использования/скриншоты (если релевантно)

## Сообщения об ошибках (issues)

- Описывайте баги максимально подробно: шаги воспроизведения, ожидаемое и фактическое поведение, окружение (OS, версия Go).
- Для feature request — опишите мотивацию и предполагаемый сценарий использования.

## Кодекс поведения

- Уважайте других участников.
- Не публикуйте вредоносный или оскорбительный контент.

## Автоматизация

- Все PR проходят автоматическую проверку (CI, линтинг, тесты, покрытие).
- Не сливайте PR без review, если вы не core-разработчик.

## Контакты

Вопросы и предложения — через issues или discussions на GitHub.

---

# Contributing to Simplex

Thank you for your interest in contributing to the Simplex project!

## How to contribute

1. **Fork the repository** and create a separate branch for your feature or fix:
   ```bash
   git checkout -b feature/your-feature
   ```
2. **Write clean, readable code**. Follow Go style and best practices.
3. **Cover your changes with tests**. Coverage for new modules should be at least 80%.
4. **Run linter and tests**:
   ```bash
   golangci-lint run
   go test ./...
   ```
5. **Write meaningful commit messages** (in English or Russian):
   - feat: new feature
   - fix: bug fix
   - docs: documentation
   - test: tests
   - refactor: refactoring
6. **Open a Pull Request**:
   - Briefly describe your changes
   - Reference related issues (if any)
   - Attach usage examples/screenshots (if relevant)

## Reporting issues

- Describe bugs in detail: steps to reproduce, expected and actual behavior, environment (OS, Go version).
- For feature requests — describe motivation and intended use case.

## Code of Conduct

- Be respectful to other contributors.
- Do not post malicious or offensive content.

## Automation

- All PRs are checked automatically (CI, lint, tests, coverage).
- Do not merge PRs without review unless you are a core developer.

## Contacts

Questions and suggestions — via issues or GitHub discussions. 