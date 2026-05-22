const fs = require('fs');
const p = require('path');

const srcDir = p.join(__dirname, '..', 'static');
const dstDir = p.join(__dirname, '..', 'dist', 'static');

fs.mkdirSync(dstDir, { recursive: true });

const icons = [
  'tab-home.svg', 'tab-home-active.svg',
  'tab-timeline.svg', 'tab-timeline-active.svg',
  'tab-pkg.svg', 'tab-pkg-active.svg',
  'tab-ai.svg', 'tab-ai-active.svg'
];

for (const f of icons) {
  const src = p.join(srcDir, f);
  const dst = p.join(dstDir, f);
  if (fs.existsSync(src)) {
    fs.copyFileSync(src, dst);
    console.log('  copied: static/' + f + ' -> dist/static/' + f);
  } else {
    console.log('  skip (not found): static/' + f);
  }
}
