#!/usr/bin/env node

const { chromium } = require('playwright');

async function scrapeTrueCaller(phoneNumber) {
    const browser = await chromium.launch({
        headless: true,
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });

    try {
        const context = await browser.newContext({
            userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
            viewport: { width: 1280, height: 720 }
        });

        const page = await context.newPage();

        // Try Kenya-specific URL first
        const url = `https://www.truecaller.com/search/ke/${phoneNumber}`;

        await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });

        // Wait for page to load
        await page.waitForTimeout(5000);

        // Get page content for debugging
        const content = await page.content();

        // Check if we need to login
        if (content.includes('Log in') || content.includes('Sign in') || content.includes('login-button')) {
            // TrueCaller requires login, try to extract any visible info
            const pageText = await page.evaluate(() => document.body.innerText);

            // Sometimes TrueCaller shows partial info even without login
            const namePatterns = [
                /Name:\s*([^\n]+)/i,
                /Owner:\s*([^\n]+)/i,
                /Registered to:\s*([^\n]+)/i
            ];

            for (const pattern of namePatterns) {
                const match = pageText.match(pattern);
                if (match && match[1]) {
                    const name = match[1].trim();
                    if (name.length > 2 && name.length < 100 && !name.includes('Unknown')) {
                        console.log(JSON.stringify({ success: true, name: name, source: 'truecaller_partial' }));
                        await browser.close();
                        return;
                    }
                }
            }
        }

        // Try to extract the name from various selectors
        const selectors = [
            '[data-testid="profile-name"]',
            '.profile-name',
            'h1[class*="name"]',
            'div[class*="profile"] h1',
            'span[class*="name"]',
            'div.name',
            'h1'
        ];

        for (const selector of selectors) {
            try {
                const elements = await page.$$(selector);
                for (const element of elements) {
                    const text = await element.textContent();
                    const name = text.trim();

                    // Validate the name
                    if (name &&
                        name !== 'Unknown' &&
                        !name.toLowerCase().includes('truecaller') &&
                        !name.toLowerCase().includes('search') &&
                        !name.toLowerCase().includes('lookup') &&
                        !name.toLowerCase().includes('reverse') &&
                        !name.toLowerCase().includes('phone') &&
                        name.length > 2 &&
                        name.length < 100) {
                        console.log(JSON.stringify({ success: true, name: name, source: 'truecaller_web' }));
                        await browser.close();
                        return;
                    }
                }
            } catch (e) {
                // Continue to next selector
            }
        }

        // Take screenshot for debugging (save to tmp)
        await page.screenshot({ path: '/tmp/truecaller_debug.png' });

        console.log(JSON.stringify({ success: false, error: 'Name not found - TrueCaller may require login', debug: 'Screenshot saved to /tmp/truecaller_debug.png' }));

    } catch (error) {
        console.log(JSON.stringify({ success: false, error: error.message }));
    } finally {
        await browser.close();
    }
}

const phoneNumber = process.argv[2];
if (!phoneNumber) {
    console.log(JSON.stringify({ success: false, error: 'No phone number provided' }));
    process.exit(1);
}

scrapeTrueCaller(phoneNumber);
