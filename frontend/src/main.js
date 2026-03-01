// This module handles the main javascript functionality for the frontend
import {
    initUI,
    getMode,
    getlink,
    showLoading,
    hideLoading,
    showEducationalResult,
    showPrankResult,
    showError,
} from './ui.js';

// Initialize UI on load
initUI();

// Global Variables
let exclude = [];

// Submit Handler
document.getElementById('submitBtn').addEventListener('click', async () => {
    const link = getlink();
    const mode = getMode();

    if (!link) {
        showError('Please enter a URL or domain before analyzing.');
        return;
    }

    showLoading();

    try {
        const response = await fetch('/api/v1/links', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ link, mode, exclude }),
        });

        const data = await response.json();
        if (!response.ok) throw new Error(data.message);

        if (mode === 'educational') {
            addToExclude(data.technique);
            resetExclude();
            showEducationalResult({
                originalLink: data.link,
                fakeLink: data.fake_link,
                technique: data.technique,
                explanation: data.explanation,
            });
        } else {
            showPrankResult(data.fake_link);
        }
        hideLoading();
    } catch (err) {
        showError(
            err.message || 'Failed to reach the server. Check your connection.',
        );
    }
});

/**
 * This function resets the exclude array if it surpasses 8 members
 */
function resetExclude() {
    if (exclude.length >= 16) {
        exclude = [];
    }
}

/**
 * This function adds a new technique string to the exclude array
 */
function addToExclude(technique) {
    exclude.push(technique);
}
