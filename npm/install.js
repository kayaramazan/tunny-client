#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const version = require('./package.json').version;
const platform = process.platform;
const arch = process.arch;

// Map Node.js platform/arch to GoReleaser naming (with title case)
const platformMap = {
  darwin: 'Darwin',
  linux: 'Linux',
  win32: 'Windows'
};

const archMap = {
  x64: 'x86_64',  // GoReleaser uses x86_64 instead of amd64
  arm64: 'arm64'
};

const goreleaserOS = platformMap[platform];
const goreleaserArch = archMap[arch];

if (!goreleaserOS || !goreleaserArch) {
  console.error(`Unsupported platform: ${platform} ${arch}`);
  process.exit(1);
}

const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `tunny${ext}`;
const binDir = path.join(__dirname, 'bin');
const binaryPath = path.join(binDir, binaryName);

// GitHub release URL (matches GoReleaser's name_template)
const downloadUrl = `https://github.com/kayaramazan/tunny-client/releases/download/v${version}/tunny_${version}_${goreleaserOS}_${goreleaserArch}.tar.gz`;

console.log(`📦 Downloading Tunny ${version} for ${platform}/${arch}...`);
console.log(`   ${downloadUrl}`);

// Create bin directory
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Download and extract
const tarPath = path.join(binDir, 'tunny.tar.gz');
const file = fs.createWriteStream(tarPath);

https.get(downloadUrl, (response) => {
  // Check for error status codes
  if (response.statusCode === 404) {
    console.error('❌ Release not found!');
    console.error(`   Version ${version} is not available on GitHub.`);
    console.error(`   Please check: ${downloadUrl}`);
    process.exit(1);
  }
  
  if (response.statusCode !== 200 && response.statusCode !== 301 && response.statusCode !== 302) {
    console.error(`❌ Download failed with status code: ${response.statusCode}`);
    process.exit(1);
  }
  
  if (response.statusCode === 302 || response.statusCode === 301) {
    // Follow redirect
    https.get(response.headers.location, (redirectResponse) => {
      if (redirectResponse.statusCode !== 200) {
        console.error(`❌ Download failed with status code: ${redirectResponse.statusCode}`);
        process.exit(1);
      }
      redirectResponse.pipe(file);
      file.on('finish', () => {
        file.close(() => extractAndCleanup());
      });
    }).on('error', (err) => {
      if (fs.existsSync(tarPath)) fs.unlinkSync(tarPath);
      console.error('❌ Download failed:', err.message);
      process.exit(1);
    });
  } else {
    response.pipe(file);
    file.on('finish', () => {
      file.close(() => extractAndCleanup());
    });
  }
}).on('error', (err) => {
  if (fs.existsSync(tarPath)) fs.unlinkSync(tarPath);
  console.error('❌ Download failed:', err.message);
  process.exit(1);
});

function extractAndCleanup() {
  try {
    // Check if file exists and is not empty
    const stats = fs.statSync(tarPath);
    if (stats.size === 0) {
      console.error('❌ Downloaded file is empty!');
      console.error('   The release might not exist on GitHub.');
      fs.unlinkSync(tarPath);
      process.exit(1);
    }
    
    // Check if file is a valid gzip (magic bytes: 1f 8b)
    const buffer = Buffer.alloc(2);
    const fd = fs.openSync(tarPath, 'r');
    fs.readSync(fd, buffer, 0, 2, 0);
    fs.closeSync(fd);
    
    if (buffer[0] !== 0x1f || buffer[1] !== 0x8b) {
      console.error('❌ Downloaded file is not a valid gzip archive!');
      console.error('   The release might not exist or the file is corrupted.');
      console.error(`   File size: ${stats.size} bytes`);
      fs.unlinkSync(tarPath);
      process.exit(1);
    }
    
    console.log('📂 Extracting...');
    
    // Extract to temporary directory
    const tempDir = path.join(binDir, 'temp');
    if (!fs.existsSync(tempDir)) {
      fs.mkdirSync(tempDir, { recursive: true });
    }
    
    if (platform === 'win32') {
      execSync(`tar -xzf "${tarPath}" -C "${tempDir}"`, { stdio: 'inherit' });
    } else {
      execSync(`tar -xzf "${tarPath}" -C "${tempDir}"`, { stdio: 'inherit' });
    }
    
    // Find the binary in the extracted files
    const files = fs.readdirSync(tempDir);
    let binaryFound = false;
    
    for (const file of files) {
      const filePath = path.join(tempDir, file);
      const stat = fs.statSync(filePath);
      
      if (stat.isDirectory()) {
        // Check inside the directory for the binary
        const innerFiles = fs.readdirSync(filePath);
        for (const innerFile of innerFiles) {
          if (innerFile === binaryName) {
            fs.renameSync(path.join(filePath, innerFile), binaryPath);
            binaryFound = true;
            break;
          }
        }
      } else if (file === binaryName) {
        // Binary is directly in temp dir
        fs.renameSync(filePath, binaryPath);
        binaryFound = true;
      }
      
      if (binaryFound) break;
    }
    
    if (!binaryFound) {
      console.error('❌ Binary not found in the archive!');
      process.exit(1);
    }
    
    // Make executable on Unix
    if (platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }
    
    // Cleanup
    fs.unlinkSync(tarPath);
    fs.rmSync(tempDir, { recursive: true, force: true });
    
    console.log('✅ Tunny installed successfully!');
    console.log(`   Run: npx tunny --help`);
  } catch (err) {
    console.error('❌ Extraction failed:', err.message);
    if (fs.existsSync(tarPath)) {
      fs.unlinkSync(tarPath);
    }
    process.exit(1);
  }
}

