#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const version = require('./package.json').version;
const platform = process.platform;
const arch = process.arch;

// Map Node.js platform/arch to Go GOOS/GOARCH
const platformMap = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64'
};

const goos = platformMap[platform];
const goarch = archMap[arch];

if (!goos || !goarch) {
  console.error(`Unsupported platform: ${platform} ${arch}`);
  process.exit(1);
}

const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `tunny${ext}`;
const binDir = path.join(__dirname, 'bin');
const binaryPath = path.join(binDir, binaryName);

// GitHub release URL
const downloadUrl = `https://github.com/yourname/tunny/releases/download/v${version}/tunny_${version}_${goos}_${goarch}.tar.gz`;

console.log(`📦 Downloading Tunny ${version} for ${goos}/${goarch}...`);
console.log(`   ${downloadUrl}`);

// Create bin directory
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Download and extract
const tarPath = path.join(binDir, 'tunny.tar.gz');
const file = fs.createWriteStream(tarPath);

https.get(downloadUrl, (response) => {
  if (response.statusCode === 302 || response.statusCode === 301) {
    // Follow redirect
    https.get(response.headers.location, (redirectResponse) => {
      redirectResponse.pipe(file);
      file.on('finish', () => {
        file.close(() => extractAndCleanup());
      });
    });
  } else {
    response.pipe(file);
    file.on('finish', () => {
      file.close(() => extractAndCleanup());
    });
  }
}).on('error', (err) => {
  fs.unlinkSync(tarPath);
  console.error('❌ Download failed:', err.message);
  process.exit(1);
});

function extractAndCleanup() {
  try {
    console.log('📂 Extracting...');
    
    if (platform === 'win32') {
      // Windows: use tar.exe or 7zip
      execSync(`tar -xzf "${tarPath}" -C "${binDir}"`, { stdio: 'inherit' });
    } else {
      // Unix: use tar
      execSync(`tar -xzf "${tarPath}" -C "${binDir}"`, { stdio: 'inherit' });
    }
    
    // Make executable on Unix
    if (platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }
    
    // Cleanup
    fs.unlinkSync(tarPath);
    
    console.log('✅ Tunny installed successfully!');
    console.log(`   Run: npx tunny --help`);
  } catch (err) {
    console.error('❌ Extraction failed:', err.message);
    process.exit(1);
  }
}

