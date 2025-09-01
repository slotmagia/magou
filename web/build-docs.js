#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

console.log('🚀 构建 HotGo 2.0 文档...\n');

const docsDir = path.join(__dirname, 'docs');
const distDir = path.join(docsDir, 'dist');

// 检查文档目录
if (!fs.existsSync(docsDir)) {
  console.error('❌ docs 目录不存在');
  process.exit(1);
}

console.log('📂 文档目录结构:');
function listFiles(dir, indent = '') {
  const items = fs.readdirSync(dir);
  items.forEach(item => {
    const fullPath = path.join(dir, item);
    const stat = fs.statSync(fullPath);
    if (stat.isDirectory()) {
      console.log(`${indent}📁 ${item}/`);
      listFiles(fullPath, indent + '  ');
    } else {
      console.log(`${indent}📄 ${item}`);
    }
  });
}

listFiles(docsDir);

console.log('\n📋 文档文件检查:');

// 检查必要的文件
const requiredFiles = [
  'index.html',
  'README.md',
  '_sidebar.md',
  '_navbar.md'
];

let allFilesExist = true;
requiredFiles.forEach(file => {
  const filePath = path.join(docsDir, file);
  if (fs.existsSync(filePath)) {
    console.log(`✅ ${file}`);
  } else {
    console.log(`❌ ${file} - 文件缺失`);
    allFilesExist = false;
  }
});

if (!allFilesExist) {
  console.error('\n❌ 缺少必要的文档文件，请检查文档结构');
  process.exit(1);
}

console.log('\n🌐 启动文档服务器...');
console.log('📍 访问地址: http://localhost:3001');
console.log('🔄 使用 Ctrl+C 停止服务\n');

try {
  // 检查是否安装了 docsify-cli
  try {
    execSync('docsify --version', { stdio: 'ignore' });
  } catch (error) {
    console.log('📦 安装 docsify-cli...');
    execSync('npm install -g docsify-cli', { stdio: 'inherit' });
  }

  // 启动文档服务器
  process.chdir(docsDir);
  execSync('docsify serve . --port 3001', { stdio: 'inherit' });
} catch (error) {
  console.error('❌ 启动文档服务器失败:', error.message);
  console.log('\n🔧 手动启动方法:');
  console.log('1. 安装 docsify-cli: npm install -g docsify-cli');
  console.log('2. 进入文档目录: cd docs');
  console.log('3. 启动服务器: docsify serve . --port 3001');
  process.exit(1);
}






