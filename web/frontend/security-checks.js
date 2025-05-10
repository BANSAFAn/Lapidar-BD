/**
 * Скрипт для проверки безопасности фронтенд-кода
 * Используется в CI/CD пайплайне для выявления уязвимостей
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Функция для проверки зависимостей на уязвимости
function checkDependencies() {
  try {
    console.log('Проверка зависимостей на уязвимости...');
    const output = execSync('npm audit --json', { encoding: 'utf8' });
    const auditResult = JSON.parse(output);
    
    if (auditResult.metadata.vulnerabilities.high > 0 || auditResult.metadata.vulnerabilities.critical > 0) {
      console.error('Обнаружены критические уязвимости в зависимостях!');
      console.error(`Критические: ${auditResult.metadata.vulnerabilities.critical}, Высокие: ${auditResult.metadata.vulnerabilities.high}`);
      process.exit(1);
    } else {
      console.log('Проверка зависимостей успешно пройдена.');
    }
  } catch (error) {
    console.error('Ошибка при проверке зависимостей:', error.message);
    process.exit(1);
  }
}

// Функция для проверки исходного кода на наличие потенциальных уязвимостей
function checkSourceCode() {
  console.log('Проверка исходного кода на потенциальные уязвимости...');
  
  const dangerousPatterns = [
    { pattern: /eval\(/, description: 'Использование eval() может привести к XSS-атакам' },
    { pattern: /document\.write\(/, description: 'document.write() может привести к XSS-уязвимостям' },
    { pattern: /innerHTML\s*=/, description: 'Прямое присваивание innerHTML может привести к XSS-уязвимостям' },
    { pattern: /dangerouslySetInnerHTML/, description: 'Проверьте использование dangerouslySetInnerHTML на безопасность' },
    { pattern: /localStorage\.setItem\(.*password/, description: 'Хранение паролей в localStorage небезопасно' },
    { pattern: /sessionStorage\.setItem\(.*password/, description: 'Хранение паролей в sessionStorage небезопасно' },
    { pattern: /console\.log\(.*password/, description: 'Логирование паролей в консоль небезопасно' },
  ];
  
  let foundIssues = false;
  
  function scanDirectory(directory) {
    const files = fs.readdirSync(directory);
    
    for (const file of files) {
      const filePath = path.join(directory, file);
      const stats = fs.statSync(filePath);
      
      if (stats.isDirectory()) {
        // Пропускаем node_modules и другие служебные директории
        if (file !== 'node_modules' && file !== 'build' && file !== 'dist') {
          scanDirectory(filePath);
        }
      } else if (stats.isFile() && /\.(js|jsx|ts|tsx)$/.test(file)) {
        // Проверяем только JavaScript/TypeScript файлы
        const content = fs.readFileSync(filePath, 'utf8');
        
        for (const { pattern, description } of dangerousPatterns) {
          if (pattern.test(content)) {
            console.error(`Потенциальная уязвимость в ${filePath}: ${description}`);
            foundIssues = true;
          }
        }
      }
    }
  }
  
  try {
    scanDirectory(path.resolve(__dirname, 'src'));
    
    if (foundIssues) {
      console.error('Обнаружены потенциальные проблемы безопасности в исходном коде!');
      process.exit(1);
    } else {
      console.log('Проверка исходного кода успешно пройдена.');
    }
  } catch (error) {
    console.error('Ошибка при проверке исходного кода:', error.message);
    process.exit(1);
  }
}

// Запуск проверок
checkDependencies();
checkSourceCode();

console.log('Все проверки безопасности успешно пройдены!');