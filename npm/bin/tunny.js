#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

const platform = process.platform;
const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `tunny${ext}`;
const binaryPath = path.join(__dirname, binaryName);

const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  env: process.env
});

child.on('exit', (code) => {
  process.exit(code);
});

