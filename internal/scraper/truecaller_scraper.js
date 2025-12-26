#!/usr/bin/env node

const { chromium } = require('playwright');

async function scrapeTrueCaller(phoneNumber) {
    const browser = await chromium.launch({
        headless: true,
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });

    