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
           