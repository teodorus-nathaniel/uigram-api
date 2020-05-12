const puppeteer = require('puppeteer');
const sharp = require('sharp');

const takeScreenshot = async (url) => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  await page.goto(url, {
    waitUntil: 'networkidle0'
  });

  const height = await page.evaluate(() => {
    const body = document.body,
      html = document.documentElement;

    return Math.max(
      body.scrollHeight,
      body.offsetHeight,
      html.clientHeight,
      html.scrollHeight,
      html.offsetHeight
    );
  });
  await page.setViewport({ width: 1920, height });
  const screenshot = await page.screenshot({
    type: 'jpeg'
  });

  const imagePath = `img/img${Math.floor(Math.random() * 100)}.jpg`;

  sharp(screenshot)
    .resize(Math.floor(1920 / 1.5), Math.floor(height / 1.5))
    .toFile(imagePath);

  await browser.close();
  return imagePath;
};

(async () => {
  try {
    const url = process.argv[2];
    const res = await takeScreenshot(url);
    console.log(res);
  } catch (error) {
    console.log(error.message);
  }
})();
