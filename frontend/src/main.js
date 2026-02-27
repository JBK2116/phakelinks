// This module handles the main javascript functionality for the frontend
import {
    initUI,
    getMode,
    getUrl,
    showLoading,
    hideLoading,
    showEducationalResult,
    showPrankResult,
    showError,
} from './ui.js';

// Initialize UI on load
initUI();

// Global Variables
const excludeTecniques = [];

// ─── Submit Handler ───────────────────────────────────────────────────────────
document.getElementById('submitBtn').addEventListener('click', async () => {
    const url = getUrl();
    const mode = getMode();

    if (!url) {
        showError('Please enter a URL or domain before analyzing.');
        return;
    }

    showLoading();

    try {
        // ── YOUR AJAX GOES HERE ──────────────────────────────────────────────────
        //
        // const response = await fetch('/api/analyze', {
        //   method: 'POST',
        //   headers: { 'Content-Type': 'application/json' },
        //   body: JSON.stringify({ url, mode })
        // });
        //
        // if (!response.ok) throw new Error(`Server error: ${response.status}`);
        // const data = await response.json();
        //
        // if (mode === 'educational') {
        //   showEduResult({
        //     originalLink: data.originalLink,
        //     fakeLink: data.fakeLink,
        //     technique: data.technique,
        //     explanation: data.explanation
        //   });
        // } else {
        //   showPrankResult(data.fakeLink);
        // }
        // ── END YOUR AJAX ────────────────────────────────────────────────────────

        // DEMO: remove this block once you wire up your backend
        await new Promise((r) => setTimeout(r, 2200));
        if (mode === 'educational') {
            showEducationalResult({
                originalLink: url,
                fakeLink:
                    url.replace(/\./g, '-').replace('https://', 'http://') +
                    '.verify-account.ru',
                technique: 'Homoglyph + Subdomain Spoofing',
                explanation:
                    'This technique replaces visually similar characters and adds a convincing subdomain to trick users into thinking they are visiting a legitimate site. The domain registers a look-alike that exploits split-second visual scanning — most users never read past the first few characters of a URL.',
            });
        } else {
            showPrankResult(
                'http://totally-not-a-virus.ru/you-won-free-iphone/claim?token=xX_' +
                    Math.random().toString(36).slice(2) +
                    '_Xx&ref=yourcredit&redirect=definitely-safe.biz',
            );
        }
    } catch (err) {
        showError(
            err.message || 'Failed to reach the server. Check your connection.',
        );
    }
});
